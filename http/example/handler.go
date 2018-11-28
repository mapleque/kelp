package example

import (
	"../../http"
)

type TodoCreateParam struct {
	Title     string `json:"title" valid:"(0,128],message=titil不合法"`
	Content   string `json:"content" valid:"optional,(0,1024),message=content不合法"`
	AlertTime string `json:"alert_time" valid:"optional,@datetime,message=alert_time不合法"`
}

func (this *Server) TodoCreate(in *TodoCreateParam) *http.Status {
	if err := this.todoList.Create(in.Title, in.Content, in.AlertTime); err != nil {
		return http.ErrorStatus(STATUS_ERROR_INTERNAL, err)
	}
	return nil
}

type TodoUpdateParam struct {
	Id        int    `json:"id" valid:"[0,),message=id不合法"`
	Title     string `json:"title" valid:"(0,128],message=titil不合法"`
	Content   string `json:"content" valid:"optional,(0,1024),message=content不合法"`
	AlertTime string `json:"alert_time" valid:"optional,@datetime,message=alert_time不合法"`
}

func (this *Server) TodoUpdate(in *TodoUpdateParam) *http.Status {
	if err := this.todoList.Update(in.Id, in.Title, in.Content, in.AlertTime); err != nil {
		return http.ErrorStatus(STATUS_ERROR_INTERNAL, err)
	}
	return nil
}

type TodoDeleteParam struct {
	Id int `json:"id" valid:"[0,),message=id不合法"`
}

func (this *Server) TodoDelete(in *TodoDeleteParam) *http.Status {
	if err := this.todoList.Delete(in.Id); err != nil {
		return http.ErrorStatus(STATUS_ERROR_INTERNAL, err)
	}
	return nil
}

type TodoRetrieveResponse struct {
	List []*Todo `json:"list"`
}

func (this *Server) TodoRetrieve(in interface{}, out *TodoRetrieveResponse) *http.Status {
	if list, err := this.todoList.Retrieve(); err != nil {
		return http.ErrorStatus(STATUS_ERROR_INTERNAL, err)
	} else {
		out.List = list
	}
	return nil
}
