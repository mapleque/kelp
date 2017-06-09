package queue

import (
	"sync"

	"github.com/kelp/log"
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

var queues map[string]*Queue

func init() {
	queues = make(map[string]*Queue)
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
	queues[name] = q
	return q
}

func GetQueue(name string) (*Queue, bool) {
	queue, ok := queues[name]
	return queue, ok
}

func GetInfo() map[string]interface{} {
	ret := make(map[string]interface{})
	for key, queue := range queues {
		ret[key] = queue.GetInfo()
	}
	return ret
}

func (queue *Queue) GetInfo() map[string]interface{} {
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

func (queue *Queue) Push(name string, flag interface{}, data interface{}) *QueueItem {
	item := &QueueItem{
		Sequence: -1,
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
		queue.mux.Lock()
		defer queue.mux.Unlock()
		qItem.Sequence = queue.sequence + 1
		queue.channel <- qItem
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
		queue.stock -= 1
		log.Info("[queue]", "pop queue", qItem)
		return qItem
	}
	log.Error("[queue]", "error queue is nil")
	return nil
}
