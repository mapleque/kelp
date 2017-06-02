package queue

import (
	"github.com/kelp/log"
)

// consumer interface
type Consumer interface {
	Pop(q *Queue, taskId string) // to consume
}

type ConsumerWrapper struct {
	Name     string
	Consumer Consumer
	Queue    *Queue
	Pause    bool
}

// a consumer map
var consumers map[string]*ConsumerWrapper

func init() {
	consumers = make(map[string]*ConsumerWrapper)
}

// regist a consumer to consumer map
func RegistConsumer(name string, c Consumer, q *Queue) {
	consumers[name] = &ConsumerWrapper{name, c, q, false}
}

func GetConsumerInfo() map[string]*ConsumerWrapper {
	return consumers
}

// run all consumer in comsumer map
func runConsumer() {
	log.Info("consumer starting ...")
	done := make(chan bool, 1)
	for _, consumer := range consumers {
		go func(c *ConsumerWrapper) {
			log.Info("consumer runing", c.Name)
			for {
				if c.Pause {
					log.Info("consumer pause", c.Name)
					break
				}
				c.Consumer.Pop(c.Queue, c.Name)
			}
		}(consumer)
	}
	<-done
}

func PauseConsumer(task string) {
	(consumers[task]).Pause = true
}
