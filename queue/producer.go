package queue

import (
	"github.com/kelp/log"
)

// producer interface
// each producer should have a *Queue property
// @see UserProducer in producer/user.go
type Producer interface {
	Push(q *Queue, taskId string) // fetch data and push into queue
}

type ProducerWrapper struct {
	Name     string
	Producer Producer
	Queue    *Queue
}

// a producer map
var producers map[string]*ProducerWrapper

func init() {
	producers = make(map[string]*ProducerWrapper)
}

// regist a producer to producer map
func RegistProducer(name string, p Producer, q *Queue) {
	producers[name] = &ProducerWrapper{name, p, q}
}

func GetProducerInfo() map[string]*ProducerWrapper {
	return producers
}

// run all producer in producer map
func runProducer() {
	log.Info("producer starting ...")
	done := make(chan bool, 1)
	for _, producer := range producers {
		go func(p *ProducerWrapper) {
			log.Info("producer runing", p.Name)
			for {
				p.Producer.Push(p.Queue, p.Name)
			}
		}(producer)
	}
	<-done
}
