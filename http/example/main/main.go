package main

import (
	service "../../example"
)

func main() {
	ts := service.NewTodoList()
	ss := service.NewServer(ts)
	ss.Init("0.0.0.0:9999")
	ss.Doc("./api.md")
	ss.Run()
}
