package queue

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

type ConsumerContainer struct {
	consumers map[string]*ConsumerWrapper
}

var cc *ConsumerContainer

func init() {
	cc = &ConsumerContainer{consumers: make(map[string]*ConsumerWrapper)}
}

func GetConsumerContainer() *ConsumerContainer {
	return cc
}

// regist a consumer to consumer map
func (q *Queue) RegistConsumer(name string, c Consumer) {
	cc.consumers[name] = &ConsumerWrapper{name, c, q, false}
}

// implement monitor.Observable
func (ccp *ConsumerContainer) GetInfo() interface{} {
	ret := make(map[string]interface{})
	for key, cw := range ccp.consumers {
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
	for _, consumer := range cc.consumers {
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
	(cc.consumers[task]).pause = true
}
