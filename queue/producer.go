package queue

import (
	"github.com/kelp/log"
)

// producer interface
type Producer interface {
	Push(q *Queue, taskId string) // fetch data and push into queue
}

type ProducerWrapper struct {
	name     string
	producer Producer
	queue    *Queue
}

// a producer map
var producers map[string]*ProducerWrapper

func init() {
	producers = make(map[string]*ProducerWrapper)
}

// regist a producer to producer map
func (q *Queue) RegistProducer(name string, p Producer) {
	producers[name] = &ProducerWrapper{name, p, q}
}

func GetProducerInfo() map[string]interface{} {
	ret := make(map[string]interface{})
	for key, pw := range producers {
		ret[key] = map[string]interface{}{
			"name":  pw.name,
			"queue": pw.queue.GetInfo()}
	}
	return ret
}

// run all producer in producer map
func runProducer() {
	log.Info("producer starting ...")
	done := make(chan bool, 1)
	for _, producer := range producers {
		go func(p *ProducerWrapper) {
			log.Info("producer runing", p.name)
			for {
				p.producer.Push(p.queue, p.name)
			}
		}(producer)
	}
	<-done
}
