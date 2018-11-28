package example

import (
	"fmt"
)

type Todo struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	AlertTime string `json:"alert_time"`
}

type TodoList struct {
	lastId int
	data   map[int]*Todo
}

func NewTodoList() *TodoList {
	return &TodoList{data: map[int]*Todo{}}
}

func (this *TodoList) Create(title, content, alertTime string) error {
	this.lastId += 1
	id := this.lastId
	if _, exist := this.data[id]; exist {
		return fmt.Errorf("duplicate id %s", id)
	}
	this.data[id] = &Todo{
		id,
		title,
		content,
		alertTime,
	}
	return nil
}

func (this *TodoList) Update(id int, title, content, alertTime string) error {
	if _, exist := this.data[id]; !exist {
		return fmt.Errorf("id is not exist %s", id)
	}
	this.data[id] = &Todo{
		id,
		title,
		content,
		alertTime,
	}
	return nil
}

func (this *TodoList) Delete(id int) error {
	if _, exist := this.data[id]; !exist {
		return fmt.Errorf("id is not exist %s", id)
	}
	delete(this.data, id)
	return nil
}

func (this *TodoList) Retrieve() ([]*Todo, error) {
	list := []*Todo{}
	for _, item := range this.data {
		list = append(list, item)
	}
	return list, nil
}
