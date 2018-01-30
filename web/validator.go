package web

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// func(fieldValue, parentValue, rootValue) pass
type ValidFunc func(reflect.Value, reflect.Value, reflect.Value) bool

var funcMap map[string]ValidFunc
var messageMap map[string]string

func init() {
	funcMap = make(map[string]ValidFunc)
	messageMap = make(map[string]string)
	funcMap["required"] = required
}

func RegisterValidFunc(key string, f ValidFunc) {
	funcMap[key] = f
}

func RegisterValidFuncWithMessage(key string, f ValidFunc, message string) {
	funcMap[key] = f
	messageMap[key] = message
}

// current field is not the default static value
func required(field reflect.Value, parent reflect.Value, root reflect.Value) bool {
	return !isEmpty(field)
}

func isEmpty(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return field.IsNil()
	default:
		return field.IsValid() && field.Interface() == reflect.Zero(field.Type()).Interface()
	}
}

// Only valid struct
// tag: valid
//		reg mode : /<regular expression>/
//		func mode : [a-z_]+
func Valid(dest interface{}) error {
	rootType := reflect.TypeOf(dest)
	if rootType.Kind() != reflect.Ptr {
		return fmt.Errorf("valid failed, dest should be a ptr but %s", rootType.Kind())
	}
	rootValue := reflect.ValueOf(dest).Elem()
	return valid(rootValue, rootValue)
}

const (
	_REG_MODE      = "/.*/"
	_RANGE_MODE    = "[\\[\\(](-{0,1}\\d+):(-{0,1}\\d+)[\\]\\)]"
	_FUNC_MODE     = "[a-zA-Z_]+[0-9a-zA-Z_]*"
	_OPTIONAL_MODE = "optional"
)

func valid(dest reflect.Value, root reflect.Value) error {
	// test each field
	for i := 0; i < dest.NumField(); i++ {
		field := dest.Type().Field(i)
		fieldValue := dest.Field(i)
		fieldName := field.Name
		fieldTags, exist := field.Tag.Lookup("valid")
		if exist {
			tags := strings.Split(fieldTags, ",")
			for _, fieldTag := range tags {
				// if optional and empty value, pass all
				if match, _ := regexp.MatchString(_OPTIONAL_MODE, fieldTag); match {
					if isEmpty(fieldValue) {
						return nil
					}
				}
			}
			for _, fieldTag := range tags {
				if match, _ := regexp.MatchString(_OPTIONAL_MODE, fieldTag); match {
					continue
				} else if match, _ := regexp.MatchString(_REG_MODE, fieldTag); match {
					// reg mode
					if field.Type.Kind() != reflect.String {
						return fmt.Errorf("valid failed, field type should be string but %s", field.Type.Kind())
					}
					regstr := fieldTag[1 : len(fieldTag)-1]
					if ok, _ := regexp.MatchString(regstr, fieldValue.String()); !ok {
						return fmt.Errorf(
							"valid failed, regexp %s test faild on %s with value %s",
							regstr, fieldName, fieldValue)
					}
				} else if match, _ := regexp.MatchString(_RANGE_MODE, fieldTag); match {
					// range mode
					r, _ := regexp.Compile(_RANGE_MODE)
					regmat := r.FindStringSubmatch(fieldTag)
					if len(regmat) != 3 {
						return fmt.Errorf("valid failed, range mode expression %s is not correct", fieldTag)
					}
					switch field.Type.Kind() {
					case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
						min, _ := strconv.ParseInt(regmat[1], 10, 64)
						max, _ := strconv.ParseInt(regmat[2], 10, 64)
						val := fieldValue.Int()
						if (min == val && fieldTag[0] != '[') || min > val {
							return fmt.Errorf("valid failed, field %s's val %d less then %d", fieldName, val, min)
						}
						if (max == val && fieldTag[len(fieldTag)-1] != ']') || max < val {
							return fmt.Errorf("valid failed, field %s's val %d great then %d", fieldName, val, max)
						}
					case reflect.Float32, reflect.Float64:
						min, _ := strconv.ParseFloat(regmat[1], 64)
						max, _ := strconv.ParseFloat(regmat[2], 64)
						val := fieldValue.Float()
						if (min == val && fieldTag[0] != '[') || min > val {
							return fmt.Errorf("valid failed, field %s's val %d less then %d", fieldName, val, min)
						}
						if (max == val && fieldTag[len(fieldTag)-1] != ']') || max < val {
							return fmt.Errorf("valid failed, field %s's val %d great then %d", fieldName, val, max)
						}
					case reflect.String:
						min, _ := strconv.Atoi(regmat[1])
						max, _ := strconv.Atoi(regmat[2])
						str := fieldValue.String()
						val := len(str)
						if (min == val && fieldTag[0] != '[') || min > val {
							return fmt.Errorf(
								"valid failed, field %s's val %s(%d) less then %d",
								fieldName, str, val, min)
						}
						if (max == val && fieldTag[len(fieldTag)-1] != ']') || max < val {
							return fmt.Errorf(
								"valid failed, field %s's val %s(%d) great then %d",
								fieldName, str, val, max)
						}
					default:
						return fmt.Errorf(
							"valid failed, range mode not support on field %s with type %s",
							fieldName, field.Type.Kind())
					}
				} else if match, _ := regexp.MatchString(_FUNC_MODE, fieldTag); match {
					// func mode
					if f, ok := funcMap[fieldTag]; ok {
						if !f(fieldValue, dest, root) {
							if msg, exist := messageMap[fieldTag]; exist {
								return fmt.Errorf(msg)
							}
							return fmt.Errorf(
								"valid failed, func %s return false on %s with value %s",
								fieldTag, fieldName, fieldValue)
						}
					} else {
						return fmt.Errorf("valid func %s need to be register", fieldTag)
					}
				} else {
					return fmt.Errorf("invalid expression in valid tag: %s", fieldTag)
				}
			}
		}
		// recursion struct field
		if field.Type.Kind() == reflect.Struct {
			if err := valid(fieldValue, root); err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *Context) Bind(dest interface{}) error {
	return Bind(this.Body, dest)
}

func Bind(data []byte, dest interface{}) error {
	if err := json.Unmarshal(data, dest); err != nil {
		return err
	}
	return Valid(dest)
}
