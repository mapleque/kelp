package http

import (
	"encoding/json"
	"reflect"
	"regexp"
)

// ValidRegWrapper 正则表达式校验包装器
// 当正则表达式中含有reflect tag中不允许的字符时，可以使用本方法生成一个校验函数
// 例如：
//  http.RegisterValidFunc("date", http.ValidRegexpWrapper(`^\d{4}-\d{2}-\d{2} \d{2}\:\d{2}\:\d{2}$`))
func ValidRegexpWrapper(reg string) ValidFunc {
	return func(
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
		return regexp.MustCompile(reg).MatchString(str)
	}
}
