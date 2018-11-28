package example

import (
	"../../http"
)

// http服务
type Server struct {
	// 放一些需要依赖的服务实例，比如mysql.Connector, logger等。

	// http server
	server *http.Server

	// 这里实现了一个简单的基于内存存储的todolist服务
	todoList *TodoList
}

// 初始化http服务
func NewServer(todoListService *TodoList) *Server {
	return &Server{todoList: todoListService}
}

func (this *Server) Init(host string) {
	// 运行后会根据router生成相应的接口文档
	server := http.New(host)
	// 注册一个log中间件，默认能够输出请求日志和错误信息
	server.Use(http.LogHandler)
	// 注册一个异常回复中间件，当handler处理过程中panic时不至于把服务搞崩
	server.Use(http.RecoveryHandler)
	// 注册一个trace中间件，在请求Header中添加一个"Kelp-Traceid: xxx"
	server.Use(http.TraceHandler)

	this.server = server

	// 注册路由
	this.initRouter()

	this.initValidator()
}

func (this *Server) Doc(file string) {
	this.server.Doc(file)
}

// 运行http服务
func (this *Server) Run() {
	// 启动服务
	this.server.Run()
}

func (this *Server) initRouter() {
	// 创建一个todo路由组
	todo := this.server.Group("/todo").Use(this.Auth) // 这里对todo路由组注册了认证的中间件
	{
		todo.Handle("添加todo", "/create", this.TodoCreate).Comment(
			"新增一条todo数据\n" +
				"```\n" +
				`curl -H'Authorization: Basic a2VscDprZWxw' -d'{"title":"kelp example","content":"write kelp example for every package","alert_time":"2018-10-26 14:00:00"}' http://host:port/todo/create` + "\n" +
				"```\n",
		)
		todo.Handle("修改todo", "/update", this.TodoUpdate).Comment(
			"修改一条todo数据\n" +
				"```\n" +
				`curl -H'Authorization: Basic a2VscDprZWxw' -d'{"id":"1","title":"kelp example","content":"write kelp example for every package","alert_time":"2018-10-26 14:00:00"}' http://host:port/todo/update` + "\n" +
				"```\n",
		)
		todo.Handle("删除todo", "/delete", this.TodoDelete).Comment(
			"删除一条todo数据\n" +
				"```\n" +
				`curl -H'Authorization: Basic a2VscDprZWxw' -d'{"id":"1"}' http://host:port/todo/delete` + "\n" +
				"```\n",
		)
		todo.Handle("列表todo", "/retrieve", this.TodoRetrieve).Comment(
			"查看所有todo数据\n" +
				"```\n" +
				`curl -H'Authorization: Basic a2VscDprZWxw' http://host:port/todo/retrieve` + "\n" +
				"```\n",
		)
	}
}

func (this *Server) initValidator() {
	// 注册一个日期(yyyy-MM-dd HH:mm:ss)的校验函数
	http.RegisterValidFunc("datetime", http.ValidRegexpWrapper(`^\d{4}-\d{2}-\d{2} \d{2}\:\d{2}\:\d{2}$`))
}
