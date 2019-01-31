package http

import (
	"encoding/json"
	"fmt"
)

type Status struct {
	Status  int         `json:"status" comment:"请参考开发者定义的Status列表"`
	Message interface{} `json:"message" comment:"用于联调测试时参考的错误信息"`
}

var (
	STATUS_SUCCESS      = &Status{0, "成功"}
	STATUS_UNKNOW       = &Status{1, "未知错误"}
	STATUS_ERROR_DB     = &Status{2, "数据库错误"}
	STATUS_UNAUTHORIZED = &Status{401, "unauthorized"}
	STATUS_FORBIDDEN    = &Status{403, "forbidden"}
	STATUS_NOT_FOUND    = &Status{404, "not found"}
)

func StatusInvalidParam(err error) *Status {
	return &Status{3, err.Error()}
}

func StatusDiy(obj interface{}) *Status {
	msg, _ := json.Marshal(obj)
	return &Status{10, string(msg)}
}

func ErrorStatus(code int, err error) *Status {
	return &Status{code, err.Error()}
}

func JsonStatus(code int, obj interface{}) *Status {
	msg, _ := json.Marshal(obj)
	return &Status{code, string(msg)}
}

func (this *Status) Error() string {
	return fmt.Sprintf("%d:%v", this.Status, this.Message)
}
