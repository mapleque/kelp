package queue

import (
	"github.com/kelp/log"
)

// consumer interface
type Consumer interface {
	Pop(q *Queue, taskId string) // to consume
}

type ConsumerWrapper struct {
	name     string
	consumer Consumer
	queue    *Queue
	pause    bool
}

// a consumer map
var consumers map[string]*ConsumerWrapper

func init() {
	consumers = make(map[string]*ConsumerWrapper)
}

// regist a consumer to consumer map
func (q *Queue) RegistConsumer(name string, c Consumer) {
	consumers[name] = &ConsumerWrapper{name, c, q, false}
}

func GetConsumerInfo() map[string]interface{} {
	ret := make(map[string]interface{})
	for key, cw := range consumers {
		ret[key] = map[string]interface{}{
			"name":  cw.name,
			"queue": cw.queue.GetInfo(),
			"pause": cw.pause}
	}
	return ret
}

// run all consumer in comsumer map
func runConsumer() {
	log.Info("consumer starting ...")
	done := make(chan bool, 1)
	for _, consumer := range consumers {
		go func(c *ConsumerWrapper) {
			log.Info("consumer runing", c.name)
			for {
				if c.pause {
					log.Info("consumer pause", c.name)
					break
				}
				c.consumer.Pop(c.queue, c.name)
			}
		}(consumer)
	}
	<-done
}

func PauseConsumer(task string) {
	(consumers[task]).pause = true
}
