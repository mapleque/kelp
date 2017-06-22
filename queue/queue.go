package queue

import (
	"sync"
)

type Queue struct {
	mux      *sync.RWMutex
	channel  chan *QueueItem
	size     int
	name     string
	sequence int
	stock    int
	flag     interface{}
}

type QueueItem struct {
	Sequence int
	Name     string
	Flag     interface{}
	Data     interface{}
}

type QueueContainer struct {
	queues map[string]*Queue
}

var qc *QueueContainer

func init() {
	qc = &QueueContainer{queues: make(map[string]*Queue)}
}

func Run() {
	done := make(chan bool, 1)
	go runProducer()
	go runConsumer()
	<-done
}

func RegistTask(name string, queueSize int, p Producer, c Consumer) *Queue {
	q := CreateQueue(name, queueSize)
	if p != nil {
		q.RegistProducer(name, p)
	}
	if q != nil {
		q.RegistConsumer(name, c)
	}
	log.Info("regist", name, q.size)
	return q
}

func CreateQueue(name string, size int) *Queue {
	q := &Queue{
		mux:      new(sync.RWMutex),
		name:     name,
		size:     size,
		channel:  make(chan *QueueItem, size),
		sequence: 0,
		stock:    0,
		flag:     nil}
	qc.queues[name] = q
	return q
}

func GetQueueContainer() *QueueContainer {
	return qc
}

func GetQueue(name string) (*Queue, bool) {
	queue, ok := qc.queues[name]
	return queue, ok
}

// implement monitor.Observable
func (qcp *QueueContainer) GetInfo() interface{} {
	ret := make(map[string]interface{})
	for key, queue := range qcp.queues {
		ret[key] = queue.GetInfo()
	}
	return ret
}

func (queue *Queue) GetInfo() interface{} {
	return map[string]interface{}{
		"size":     queue.size,
		"name":     queue.name,
		"stock":    queue.stock,
		"sequence": queue.sequence,
		"flag":     queue.flag}
}

func (queue *Queue) GetFlag() interface{} {
	return queue.flag
}

func (queue *Queue) SetFlag(flag interface{}) {
	queue.flag = flag
}

func (queue *Queue) GetSize() int {
	return queue.size
}

func (queue *Queue) Push(name string, flag interface{}, data interface{}) *QueueItem {
	item := &QueueItem{
		Sequence: queue.sequence + 1,
		Name:     name,
		Flag:     flag,
		Data:     data}
	queue.push(item)
	log.Info("[queue]", "push queue", item)
	return item
}

func (queue *Queue) push(qItem *QueueItem) {
	if queue != nil {
		overStock := false
		if queue.stock == queue.size {
			overStock = true
			log.Warn("[queue]", "waiting ..., data over stock on ", qItem)
		}
		queue.channel <- qItem
		queue.mux.Lock()
		defer queue.mux.Unlock()
		queue.sequence += 1
		queue.stock += 1
		queue.flag = qItem.Flag
		if overStock {
			log.Warn("[queue]", "over stock recover on ", qItem)
		}
	}
}

func (queue *Queue) Pop() *QueueItem {
	if queue != nil {
		qItem := <-queue.channel
		queue.mux.Lock()
		defer queue.mux.Unlock()
		queue.stock -= 1
		log.Info("[queue]", "pop queue", qItem)
		return qItem
	}
	log.Error("[queue]", "error queue is nil")
	return nil
}
