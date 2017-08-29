package web

import (
	"fmt"
)

type ParamCheckFunc func(interface{}) bool
type ParamTransFunc func(*Context, string)

func Exist(context *Context, field string) bool {
	return context.GetParam(field) != nil
}

func OptionalCheck(
	context *Context, field string, checkFunc ParamCheckFunc, transFunc ParamTransFunc, errorCode int) bool {
	return !Exist(context, field) || Check(context, field, checkFunc, transFunc, errorCode)
}

func Check(context *Context, field string, checkFunc ParamCheckFunc, transFunc ParamTransFunc, errorCode int) bool {
	param := context.GetParam(field)
	if param == nil {
		context.Status = errorCode
		return false
	}
	if checkFunc(param) {
		if transFunc != nil {
			transFunc(context, field)
		}
		return true
	}
	context.Status = errorCode
	return false
}

func InArray(arr []interface{}) ParamCheckFunc {
	return func(param interface{}) bool {
		for _, v := range arr {
			if equal(v, param) {
				return true
			}
		}
		return false
	}
}

func IsSubSet(arr []interface{}) ParamCheckFunc {
	return func(param interface{}) bool {
		switch ret := param.(type) {
		case []interface{}:
			for _, b := range ret {
				if !InArray(arr)(b) {
					log.Debug("not in array", arr, b)
					return false
				}
			}
			return true
		default:
			log.Debug("type assert error", fmt.Sprintf("%T", ret))
		}
		return false
	}
}

func IsString(minLen, maxLen int) ParamCheckFunc {
	return func(param interface{}) bool {
		switch ret := param.(type) {
		case string:
			// TODO prevent special charactor
			strLen := len(ret)
			if minLen >= 0 && strLen < minLen {
				log.Debug(fmt.Sprintf("string len %d less then minLen %d", strLen, minLen))
				return false
			}
			if maxLen >= 0 && strLen > maxLen {
				log.Debug(fmt.Sprintf("string len %d more then maxLen %d", strLen, maxLen))
				return false
			}
			return true
		default:
			log.Debug("type assert error", fmt.Sprintf("%T", ret))
		}
		return false
	}
}

func IsNum(min, max int) ParamCheckFunc {
	return func(param interface{}) bool {
		switch ret := param.(type) {
		case int, float64:
			val := toInt(ret)
			if min >= 0 && val < min {
				log.Debug(fmt.Sprintf("num %d less then min %d", ret, min))
				return false
			}
			if max >= 0 && val > max {
				log.Debug(fmt.Sprintf("num %d more then max %d", ret, max))
				return false
			}
			return true
		default:
			log.Debug("type assert error", fmt.Sprintf("%T", ret))
		}
		return false
	}
}

func toInt(val interface{}) int {
	switch ret := val.(type) {
	case int:
		return ret
	case float64:
		return int(ret)
	}
	return 0
}

func equal(a interface{}, b interface{}) bool {
	switch at := a.(type) {
	case string:
		switch bt := b.(type) {
		case string:
			if at == bt {
				return true
			}
		default:
			log.Debug("type assert error", fmt.Sprintf("a:%T,b:%T", at, bt))
			return false
		}
	case int:
		switch bt := b.(type) {
		case int:
			if at == bt {
				return true
			}
		case float64:
			if at == int(bt) {
				return true
			}
		default:
			log.Debug("type assert error", fmt.Sprintf("a:%T,b:%T", at, bt))
			return false
		}
	case float64:
		switch bt := b.(type) {
		case float64:
			if at == bt {
				return true
			}
		case int:
			if int(at) == bt {
				return true
			}
		default:
			log.Debug("type assert error", fmt.Sprintf("a:%T,b:%T", at, bt))
			return false
		}
	default:
		return a == b
	}
	return false
}
