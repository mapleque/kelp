package util

import (
	"fmt"
	"strconv"

	"github.com/kelp/log"
)

// 类型转换，任何类型转成int
func Int(param interface{}) int {
	switch ret := param.(type) {
	case int:
		return ret
	case int64:
		return int(ret)
	case float64:
		return int(ret)
	case string:
		r, err := strconv.Atoi(ret)
		if err != nil {
			log.Error("param type change error", ret, err.Error())
		}
		return r
	case bool:
		if ret {
			return 1
		} else {
			return 0
		}
	default:
		log.Error("param type change to int error",
			ret, fmt.Sprintf("%T", ret))
		return 0
	}
}

// 类型转换，类型转换成float
func Float(param interface{}) float64 {
	switch ret := param.(type) {
	case int:
		return float64(ret)
	case int64:
		return float64(ret)
	case float64:
		return ret
	case string:
		r, err := strconv.ParseFloat(ret, 64)
		if err != nil {
			log.Error("param type change error", ret, err.Error())
		}
		return r
	case bool:
		if ret {
			return 1.0
		} else {
			return 0.0
		}
	default:
		log.Error("param type change to int error",
			ret, fmt.Sprintf("%T", ret))
		return 0.0
	}
}

// 类型转换，任何类型转成bool
func Bool(param interface{}) bool {
	switch ret := param.(type) {
	case bool:
		return ret
	case int:
		if ret > 0 {
			return true
		} else {
			return false
		}
	case string:
		switch ret {
		case "1", "true", "y", "on", "yes":
			return true
		case "0", "false", "n", "off", "no":
			return false
		default:
			log.Error("param type change to bool error", ret, "unknown type")
		}
		return false
	default:
		log.Error("param type change to bool error",
			ret, fmt.Sprintf("%T", ret))
		return false
	}
}

// 类型转换，任何类型转成string
func String(param interface{}) string {
	switch ret := param.(type) {
	case string:
		return ret
	case int:
		return strconv.Itoa(ret)
	case bool:
		if ret {
			return "1"
		} else {
			return "0"
		}
	default:
		log.Error("param type change to string error",
			ret, fmt.Sprintf("%T", ret))
		return ""
	}
}

func Map(param interface{}) map[string]interface{} {
	switch ret := param.(type) {
	case map[string]interface{}:
		return ret
	default:
		log.Error("param type change to map error",
			ret, fmt.Sprintf("%T", ret))
		return nil
	}
}

func Array(param interface{}) []interface{} {
	switch ret := param.(type) {
	case []interface{}:
		return ret
	default:
		log.Error("param type change to map error",
			ret, fmt.Sprintf("%T", ret))
		return nil
	}
}
