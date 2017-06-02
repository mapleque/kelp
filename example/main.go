package main

import (
	"flag"
	"time"

	"github.com/kelp/config"
	"github.com/kelp/crontab"
	"github.com/kelp/db"
	"github.com/kelp/log"
	"github.com/kelp/monitor"
	"github.com/kelp/queue"
)

var conf config.Configer

type SimpleImpl struct{}

func (p SimpleImpl) Push(q *queue.Queue, task string) {
	qItem := q.Push(task, 1, "item data")
	log.Info("[producer", task, "]", "push qItem", qItem)
}

func (c SimpleImpl) Pop(q *queue.Queue, task string) {
	qItem := q.Pop()
	log.Info("[consumer", task, "]", "fetch and deal", qItem)
	time.Sleep(2000000000)
}

func (c SimpleImpl) Triger(task string) {
	log.Info("[crontab ", task, "]", "do sth ...")
}

func initLog() {
	log.AddLogger(
		conf.Get("log.NAME"),
		conf.Get("log.PATH"),
		conf.Int("log.MAX_NUMBER"),
		conf.Int64("log.MAX_SIZE"),
		conf.Int("log.MAX_LEVEL"),
		conf.Int("log.MIN_LEVEL"))
}

func initDB() {
	db.AddDB(
		"kelp",
		conf.Get("db.DSN"),
		conf.Int("db.MAX_CONNECTION"),
		conf.Int("db.MAX_IDLE"))
}

func main() {
	confFile := flag.String("ini", "./config.ini", "your config file")
	flag.Parse()
	if *confFile == "" {
		panic("run with -h to find usage")
	}
	config.AddConfiger(config.INI, "config_name", *confFile)
	conf = config.Use("config_name")

	initLog()
	log.Info("init db")
	initDB()
	log.Info("start regist ...")

	impl := SimpleImpl{}
	queue.RegistTask("simple", 10, impl, impl)
	crontab.Regist("* * * * *", "simple", impl)

	// start service
	log.Info("start service ...")
	done := make(chan bool, 1)

	go queue.Run()
	go crontab.Run()
	go monitor.Run(conf.Get("monitor.HOST"))

	<-done
}
