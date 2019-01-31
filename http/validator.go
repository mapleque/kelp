package http

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ValidFunc 校验函数类型
// valid tag中使用的函数必须是这个类型的
type ValidFunc func(
	fieldType reflect.StructField,
	destSource []byte,
	root reflect.Value,
	rootSource map[string]json.RawMessage,
) bool

var (
	funcMap        map[string]ValidFunc
	validRuleCache map[string]*validRuleGroup
)

func init() {
	funcMap = make(map[string]ValidFunc)
	validRuleCache = make(map[string]*validRuleGroup)
}

// RegisterValidFunc 注册一个校验函数
// 只有注册后的函数才能够在valid tag中使用
// 如果使用了没有注册的函数，则该校验始终返回false
func RegisterValidFunc(key string, f ValidFunc) {
	funcMap[key] = f
}

// Valid 校验一个json source是否满足目标dest struct中声明的valid tag要求
func Valid(dest interface{}, source []byte) error {
	// 获取目标类型
	rootType := reflect.TypeOf(dest)
	// 如果该类型不是指针，就直接返回错误
	// 因为不是指针的实体无法被绑定数据
	if rootType.Kind() != reflect.Ptr {
		return fmt.Errorf("[this is a system error should be fixed by developer] dest should be a ptr but %s", rootType.Kind())
	}
	// 获取目标实体
	rootValue := reflect.ValueOf(dest).Elem()

	// 将source绑定到map上
	sourceValue := map[string]json.RawMessage{}
	if err := json.Unmarshal(source, &sourceValue); err != nil {
		return err
	}
	return valid(rootValue, sourceValue, rootValue, sourceValue)
}

func valid(
	dest reflect.Value,
	destSource map[string]json.RawMessage,
	root reflect.Value,
	rootSource map[string]json.RawMessage,
) error {
	// 先获取当前对象的所有属性
	for i := 0; i < dest.NumField(); i++ {
		fieldType := dest.Type().Field(i)
		fieldName := getFieldName(fieldType)
		fieldValue := dest.Field(i)
		if fieldTags, exist := fieldType.Tag.Lookup("valid"); exist {
			// 如果有校验标记，则进行校验
			if err := validField(fieldType, fieldTags, destSource, root, rootSource); err != nil {
				return err
			}
		}

		// 如果是struct，则需要继续递归，因为struct内部可能还有valid
		if fieldType.Type.Kind() == reflect.Struct {
			// 这里要判断目标数据是否存在，如果是nil，则把nil继续传递下去
			if destSource == nil {
				if err := valid(fieldValue, nil, root, rootSource); err != nil {
					return err
				}
			}
			if source, ok := destSource[fieldName]; ok {
				// 这里要判断目标数据是否存在这个字段，如果不存在，就都替换成nil
				if err := valid(fieldValue, nil, root, rootSource); err != nil {
					return err
				}
			} else {
				srcBytes, _ := source.MarshalJSON()
				src := map[string]json.RawMessage{}
				// 这里因为确定这层必须是struct，所以把source再转成map[string]interface{}
				// 如果在转换的时候出错，说明数据类型对不上，报错
				if err := json.Unmarshal(srcBytes, &src); err != nil {
					return err
				}
				// 最后递归
				if err := valid(fieldValue, src, root, rootSource); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func getFieldName(fieldType reflect.StructField) string {
	if jsonTag, exist := fieldType.Tag.Lookup("json"); exist {
		return strings.Split(jsonTag, ",")[0]
	}
	return fieldType.Name
}

func validField(
	fieldType reflect.StructField,
	fieldTags string,
	destSource map[string]json.RawMessage,
	root reflect.Value,
	rootSource map[string]json.RawMessage,
) error {
	fieldName := getFieldName(fieldType)
	if ruleGroup, err := validParse(fieldTags); err != nil {
		// 这里要出错，说明定义valid tag的有问题
		return err
	} else {
		// 先看destSource有没有，如果没有，就可以根据optional标记直接返回了
		if destSource == nil {
			if ruleGroup.optional {
				return nil
			} else {
				return ruleGroup.err(fieldName)
			}
		}
		if source, exist := destSource[fieldName]; !exist {
			if ruleGroup.optional {
				return nil
			} else {
				return ruleGroup.err(fieldName)
			}
		} else {
			// 再看是否满足其他要求，将所有rule都过一遍即可
			srcBytes, _ := source.MarshalJSON()
			for _, rule := range ruleGroup.rules {
				if !rule.valid(fieldType, srcBytes, root, rootSource) {
					return ruleGroup.err(fieldName)
				}
			}
		}
	}
	return nil
}

type validRuleGroup struct {
	rules    []rule
	optional bool
	message  string
}

func (this *validRuleGroup) err(fieldName string) error {
	if this.message != "" {
		return fmt.Errorf(this.message)
	}
	return fmt.Errorf("%s valid faild", fieldName)
}

type rule interface {
	// 返回是否通过校验
	valid(
		fieldType reflect.StructField,
		destSource []byte,
		root reflect.Value,
		rootSource map[string]json.RawMessage,
	) bool
	// 返回是哪种rule
	value() string
}

const (
	_REG_MODE          = "/.*/"
	_REG_MODE_EXEC     = "/(.*)/"
	_RANGE_MODE        = `[\[\(]-{0,1}\d*,-{0,1}\d*[\]\)]`
	_RANGE_MODE_EXEC   = `([\[\(])(-{0,1}\d*),(-{0,1}\d*)([\]\)])`
	_FUNC_MODE         = "@[a-zA-Z_]+[0-9a-zA-Z_]*"
	_OPTIONAL_MODE     = "optional"
	_MESSAGE_MODE      = `message=.*`
	_MESSAGE_MODE_EXEC = `message=(.*)`
)

func validParse(tagsField string) (*validRuleGroup, error) {
	if cacheRuleGroup, exist := validRuleCache[tagsField]; exist {
		return cacheRuleGroup, nil
	}

	ret := &validRuleGroup{
		rules:    []rule{},
		optional: false,
		message:  "",
	}
	// 如果valid tag是空字符串，那么就只验证存在性
	// 因为这里optional默认是false,rules默认是空数组
	if tagsField == "" {
		validRuleCache[tagsField] = ret
		return ret, nil
	}

	matchArr := []string{}
	for _, mode := range []string{
		_OPTIONAL_MODE,
		_MESSAGE_MODE,
		_REG_MODE,
		_RANGE_MODE,
		_FUNC_MODE,
	} {
		if m := regexp.MustCompile("("+mode+")").FindAllStringSubmatch(tagsField, -1); m != nil && len(m) > 0 {
			for _, tm := range m {
				matchArr = append(matchArr, tm[1])
			}
		}
	}

	// 遍历所有匹配到的规则，构造validRuleGroup
	for _, rule := range matchArr {
		switch {
		case rule == "":
			continue
		case regexp.MustCompile(_OPTIONAL_MODE).MatchString(rule):
			ret.optional = true
		case regexp.MustCompile(_MESSAGE_MODE).MatchString(rule):
			ret.message = regexp.MustCompile(_MESSAGE_MODE_EXEC).FindStringSubmatch(rule)[1]
		case regexp.MustCompile(_REG_MODE).MatchString(rule):
			ret.rules = append(ret.rules, newRegRule(
				regexp.MustCompile(_REG_MODE_EXEC).FindStringSubmatch(rule)[1],
			))
		case regexp.MustCompile(_RANGE_MODE).MatchString(rule):
			ele := regexp.MustCompile(_RANGE_MODE_EXEC).FindStringSubmatch(rule)
			ret.rules = append(ret.rules, newRangeRule(ele[2], ele[3], ele[1] == "[", ele[4] == "]"))
		case regexp.MustCompile(_FUNC_MODE).MatchString(rule):
			ret.rules = append(ret.rules, newFuncRule(rule[1:]))
		default:
			panic("regexp not match " + rule)
		}
	}

	validRuleCache[tagsField] = ret
	return ret, nil
}

type regRule struct {
	reg string
}

func newRegRule(reg string) *regRule {
	return &regRule{reg}
}
func (this *regRule) value() string {
	return this.reg
}

func (this *regRule) valid(
	fieldType reflect.StructField,
	destSource []byte,
	root reflect.Value,
	rootSource map[string]json.RawMessage,
) bool {
	str := string(destSource)
	val := len(str)
	// 如果是字符串，去掉两端的引号
	if val > 2 && str[0] == '"' && str[val-1] == '"' {
		str = str[1 : val-1]
	}
	return regexp.MustCompile(this.reg).MatchString(str)
}

type rangeRule struct {
	min      string
	max      string
	equalMin bool
	equalMax bool
}

func newRangeRule(min, max string, equalMin, equalMax bool) *rangeRule {
	return &rangeRule{min, max, equalMin, equalMax}
}

func (this *rangeRule) value() string {
	emin := ""
	emax := ""
	if this.equalMin {
		emin = "="
	}
	if this.equalMax {
		emax = "="
	}
	return fmt.Sprintf("%s<%sx<%s%s", this.min, emin, emax, this.max)
}

func (this *rangeRule) valid(
	fieldType reflect.StructField,
	destSource []byte,
	root reflect.Value,
	rootSource map[string]json.RawMessage,
) bool {
	switch fieldType.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		// 如果是int，就比较int值大小
		val, err := strconv.ParseInt(string(destSource), 10, 64)
		if err != nil {
			return false
		}
		switch {
		case this.min == "" && this.max == "":
			return true
		case this.min == "":
			max, _ := strconv.ParseInt(this.max, 10, 64)
			return val < max || (val == max && this.equalMax)
		case this.max == "":
			min, _ := strconv.ParseInt(this.min, 10, 64)
			return val > min || (val == min && this.equalMin)
		default:
			min, _ := strconv.ParseInt(this.min, 10, 64)
			max, _ := strconv.ParseInt(this.max, 10, 64)
			return (val > min || (val == min && this.equalMin)) &&
				(val < max || (val == max && this.equalMax))
		}
	case reflect.Float32, reflect.Float64:
		// 如果是float，就比较float值大小
		val, err := strconv.ParseFloat(string(destSource), 64)
		if err != nil {
			return false
		}
		switch {
		case this.min == "" && this.max == "":
			return true
		case this.min == "":
			max, _ := strconv.ParseFloat(this.max, 64)
			return val < max || (val == max && this.equalMax)
		case this.max == "":
			min, _ := strconv.ParseFloat(this.min, 64)
			return val > min || (val == min && this.equalMin)
		default:
			min, _ := strconv.ParseFloat(this.min, 64)
			max, _ := strconv.ParseFloat(this.max, 64)
			return (val > min || (val == min && this.equalMin)) &&
				(val < max || (val == max && this.equalMax))
		}
	case reflect.String:
		// 如果是string，就比较长度
		str := string(destSource)
		val := len(str)
		// 目标类型不对
		if str[0] != '"' || str[val-1] != '"' {
			return false
		}
		// 去掉两端的引号
		val = val - 2
		switch {
		case this.min == "" && this.max == "":
			return true
		case this.min == "":
			max, _ := strconv.Atoi(this.max)
			return val < max || (val == max && this.equalMax)
		case this.max == "":
			min, _ := strconv.Atoi(this.min)
			return val > min || (val == min && this.equalMin)
		default:
			max, _ := strconv.Atoi(this.max)
			min, _ := strconv.Atoi(this.min)
			return (val > min || (val == min && this.equalMin)) &&
				(val < max || (val == max && this.equalMax))
		}
	default:
		// 如果是其他的，就都算不通过
		return false
	}
}

type funcRule struct {
	funcName string
}

func newFuncRule(funcName string) *funcRule {
	return &funcRule{funcName}
}
func (this *funcRule) value() string {
	return this.funcName
}

func (this *funcRule) valid(
	fieldType reflect.StructField,
	destSource []byte,
	root reflect.Value,
	rootSource map[string]json.RawMessage,
) bool {
	if f, exist := funcMap[this.funcName]; !exist {
		return false
	} else {
		return f(fieldType, destSource, root, rootSource)
	}
}

// BindAndValidJson 将请求body绑定到目标类型实体上
// body的内容必须是合法的json格式
func (this *Context) BindAndValidJson(dest interface{}) error {
	return BindAndValidJson(dest, this.Body())
}

// BindAndValidJson 将data绑定到目标类型实体上
// body的内容必须是合法的json格式
func BindAndValidJson(dest interface{}, data []byte) error {
	if err := json.Unmarshal(data, dest); err != nil {
		errType := reflect.TypeOf(err).Elem()
		switch errType.Name() {
		case "UnmarshalTypeError": // 处理合法json类型非法的情况
			// 先查找有没有预定义的message
			errReal := reflect.ValueOf(err).Interface().(*json.UnmarshalTypeError)
			if errReal.Struct != "" && errReal.Field != "" {
				destType := reflect.TypeOf(dest).Elem()
				var field reflect.StructField
				var exist bool
				if destType.Name() == errReal.Struct {
					if field, exist = findFieldByName(destType, errReal.Field); exist {
						if validTag, exist := field.Tag.Lookup("valid"); exist {
							if regexp.MustCompile(_MESSAGE_MODE).MatchString(validTag) {
								message := regexp.MustCompile(_MESSAGE_MODE_EXEC).FindStringSubmatch(validTag)[1]
								return fmt.Errorf(message)
							}
						} else {
							return fmt.Errorf("invalid json %s with error %v", string(data), err)
						}
					}
				}
				// 处理嵌套情况
				for i := 0; i < destType.NumField(); i++ {
					rootField := destType.Field(i)
					if rootField.Type.Kind() == reflect.Struct {
						if field, exist = findField(rootField, errReal.Struct, errReal.Field); exist {
							if validTag, exist := field.Tag.Lookup("valid"); exist {
								if regexp.MustCompile(_MESSAGE_MODE).MatchString(validTag) {
									message := regexp.MustCompile(_MESSAGE_MODE_EXEC).FindStringSubmatch(validTag)[1]
									return fmt.Errorf(message)
								}
							} else {
								return fmt.Errorf("invalid json %s with error %v", string(data), err)
							}
						}
					}
				}
			}
			// 不需要处理一层情况

			// 没有message返回默认错误
			return fmt.Errorf("invalid json %s with error %v", string(data), err)
		// case "UnmarshalFieldError": // 处理合法json字段非法的情况
		// 	// TODO deal this
		// 	return fmt.Errorf("invalid json %s with error %v", string(data), err)
		// case "InvalidUnmarshalError": // 处理非法json的情况
		// 	// TODO deal this
		// 	return fmt.Errorf("invalid json %s with error %v", string(data), err)
		default:
			return fmt.Errorf("invalid json %s with error %v", string(data), err)
		}
	}
	return Valid(dest, data)
}

func findField(dest reflect.StructField, structName, fieldName string) (field reflect.StructField, exist bool) {
	// 广度优先
	if field, exist := findFieldByName(dest.Type, fieldName); exist {
		return field, true
	}
	// 再去递归
	for i := 0; i < dest.Type.NumField(); i++ {
		field := dest.Type.Field(i)
		// 如果是struct，就在其内部再去搜索
		if field.Type.Kind() == reflect.Struct {
			if result, exist := findField(field, structName, fieldName); exist {
				return result, exist
			}
		}
	}
	return dest, false
}

func findFieldByName(dest reflect.Type, fieldName string) (field reflect.StructField, exist bool) {
	for i := 0; i < dest.NumField(); i++ {
		field := dest.Field(i)
		if getFieldName(field) == fieldName {
			return field, true
		}
	}
	return reflect.StructField{}, false
}
