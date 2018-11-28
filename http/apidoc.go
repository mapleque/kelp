package http

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type docBuilder struct {
	apiDocs []*apiDoc
}

type apiDoc struct {
	title   string
	path    string
	comment string

	param    interface{}
	response interface{}
	status   interface{}
}

type defaultSuccessResponse struct {
	Status int         `json:"status" comment:"默认为0"`
	Data   interface{} `json:"data" comment:"默认为空字符串"`
}

func (this *Server) Doc(path string) {
	db := &docBuilder{}
	db.buildRouterDoc(this.router)
	db.output(path)
}

func (this *docBuilder) output(path string) {
	var file *os.File
	var err error
	if path == "" {
		file = os.Stdout
	} else {
		file, err = os.Create(path)
		if err != nil {
			panic(err)
		}
	}
	fmt.Fprintln(file, "# 接口文档")
	for _, apiDoc := range this.apiDocs {
		fmt.Fprintln(file)
		fmt.Fprintln(file, "##", apiDoc.title)
		if apiDoc.comment != "" {
			fmt.Fprintln(file)
			fmt.Fprintln(file, apiDoc.comment)
		}
		fmt.Fprintln(file)
		fmt.Fprintln(file, "请求路径：`", apiDoc.path, "`")
		fmt.Fprintln(file)
		fmt.Fprintln(file, outputJson("请求参数：", apiDoc.param))
		fmt.Fprintln(file, outputJson("返回数据：", apiDoc.response))
		fmt.Fprintln(file, outputJson("异常返回：", apiDoc.status))
	}
}

func outputJson(title string, jsonObj interface{}) string {
	ret := ""
	if jsonObj != nil {
		if c := json2String(generalJsonDoc(jsonObj, "", "", "")); c != "" && c != `""` {
			ret += fmt.Sprintln(title)
			ret += fmt.Sprintln("```")
			ret += fmt.Sprintln(c)
			ret += fmt.Sprintln("```")
		}
	}
	return ret
}

func (this *docBuilder) buildRouterDoc(router *Router) {
	doc := &apiDoc{
		param:    nil,
		response: &defaultSuccessResponse{},
		status:   &Status{},
	}
	if router.title != "" {
		doc.title = router.title
		doc.comment = router.comment
		doc.path = router.realPath
		doc.buildHandlerDoc(router.handlerChain)
		this.apiDocs = append(this.apiDocs, doc)
	}
	for _, r := range router.children {
		this.buildRouterDoc(r)
	}
}

func (this *apiDoc) buildHandlerDoc(handlerChain []HandlerFunc) {
	for _, handlerFunc := range handlerChain {
		handlerType := reflect.TypeOf(handlerFunc)
		if handlerType.Kind() != reflect.Func {
			panic("handler type must be func but " + handlerType.Name())
		}
		switch handlerType.NumIn() {
		case 1:
			paramType := handlerType.In(0)
			if paramType.Elem().Name() != "Context" {
				this.param = reflect.New(paramType.Elem()).Interface()
			}
			// param is in
			// response is default
		case 2, 3:
			paramType := handlerType.In(0)
			if paramType.Kind() == reflect.Ptr {
				this.param = reflect.New(paramType.Elem()).Interface()
			} else {
				this.param = reflect.New(paramType).Interface()
			}

			responseType := handlerType.In(1)
			if responseType.Kind() == reflect.Ptr {
				this.response = reflect.New(responseType.Elem()).Interface()
			} else {
				this.response = reflect.New(responseType).Interface()
			}

			// param is in
			// response is out
		default:
			// param is nil
			// response is default
		}
		if handlerType.NumOut() > 0 {
			statusType := handlerType.Out(0)
			if statusType.Kind() == reflect.Ptr {
				this.status = reflect.New(statusType.Elem()).Interface()
			} else {
				this.status = reflect.New(statusType).Interface()
			}
		}
	}
}

func json2String(dest interface{}) string {
	bytes, _ := json.MarshalIndent(dest, "", "  ")
	return string(bytes)
}

func generalJsonDoc(obj interface{}, kind, rule, comment string) interface{} {
	if obj != nil {
		objType := reflect.TypeOf(obj)
		switch objType.Kind() {
		case reflect.Ptr:
			objValue := reflect.ValueOf(obj)
			if !objValue.IsNil() {
				return generalJsonDoc(reflect.ValueOf(obj).Elem().Interface(), kind, rule, comment)
			} else {
				return generalJsonDoc(reflect.New(objType.Elem()).Interface(), kind, rule, comment)
			}
		case reflect.Struct:
			ret := map[string]interface{}{}
			for i := 0; i < objType.NumField(); i++ {
				fieldType := objType.Field(i)
				fieldName := getJsonName(fieldType)
				if fieldName != "" && fieldName != "-" {
					ret[fieldName] = generalJsonDoc(
						reflect.New(fieldType.Type).Interface(),
						fieldType.Type.Kind().String(),
						getFieldTag(fieldType, "valid"),
						getFieldTag(fieldType, "comment"),
					)
				}
			}
			return ret
		case reflect.Map:
			ret := map[string]interface{}{}
			objValue := reflect.ValueOf(obj)
			for _, key := range objValue.MapKeys() {
				if reflect.TypeOf(key).Kind() == reflect.String {
					ret[key.String()] = generalJsonDoc(objValue.MapIndex(key).Interface(), "", "", "")
				}
			}
			return ret
		case reflect.Slice:
			return []interface{}{
				generalJsonDoc(reflect.New(objType.Elem()).Interface(), "", "", ""),
			}
		default:
			// do nothing
		}
	}
	ret := ""
	if kind != "" {
		ret += fmt.Sprintf("%s", kind)
	}
	if rule != "" {
		ret += fmt.Sprintf(" |%s|", rule)
	}
	if comment != "" {
		ret += fmt.Sprintf(" // %s", comment)
	}
	return ret
}

func getJsonName(fieldType reflect.StructField) string {
	if tag, exist := fieldType.Tag.Lookup("json"); exist {
		return strings.Split(tag, ",")[0]
	}
	return ""
}

func getFieldTag(fieldType reflect.StructField, tagName string) string {
	if tag, exist := fieldType.Tag.Lookup(tagName); exist {
		return tag
	}
	return ""
}
