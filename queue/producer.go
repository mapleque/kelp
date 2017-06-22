package queue

// producer interface
type Producer interface {
	Push(q *Queue, taskId string) // fetch data and push into queue
}

type ProducerWrapper struct {
	name     string
	producer Producer
	queue    *Queue
}

type ProducerContainer struct {
	producers map[string]*ProducerWrapper
}

var pc *ProducerContainer

func init() {
	pc = &ProducerContainer{producers: make(map[string]*ProducerWrapper)}
}

func GetProducerContainer() *ProducerContainer {
	return pc
}

// regist a producer to producer map
func (q *Queue) RegistProducer(name string, p Producer) {
	pc.producers[name] = &ProducerWrapper{name, p, q}
}

// implement monitor.Observable
func (pcp *ProducerContainer) GetInfo() interface{} {
	ret := make(map[string]interface{})
	for key, pw := range pcp.producers {
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
	for _, producer := range pc.producers {
		go func(p *ProducerWrapper) {
			log.Info("producer runing", p.name)
			for {
				p.producer.Push(p.queue, p.name)
			}
		}(producer)
	}
	<-done
}
