package queue

import (
	"sync"

	"github.com/kelp/log"
)

type Queue struct {
	mux             *sync.RWMutex
	Size            int
	stock           chan *QueueItem
	CurrentSequence int
	CurrentStock    int
	CurrentPid      int
}

type QueueItem struct {
	Sequence int
	Name     string
	Pid      int
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
	log.Info("regist", name, q.Size)
	return q
}

func CreateQueue(name string, size int) *Queue {
	q := &Queue{
		mux:             new(sync.RWMutex),
		Size:            size,
		stock:           make(chan *QueueItem, size),
		CurrentSequence: 0,
		CurrentStock:    0,
		CurrentPid:      0}
	queues[name] = q
	return q
}

func GetQueue(name string) (*Queue, bool) {
	queue, ok := queues[name]
	return queue, ok
}

func GetInfo() map[string]*Queue {
	return queues
}

func (queue *Queue) Push(name string, pid int, data interface{}) *QueueItem {
	item := &QueueItem{
		Sequence: -1,
		Name:     name,
		Pid:      pid,
		Data:     data}
	queue.push(item)
	log.Info("[queue]", "push queue", item)
	return item
}

func (queue *Queue) push(qItem *QueueItem) {
	if queue != nil {
		overStock := false
		if queue.CurrentStock == queue.Size {
			overStock = true
			log.Warn("[queue]", "waiting ..., data over stock on ", qItem)
		}
		queue.mux.Lock()
		defer queue.mux.Unlock()
		qItem.Sequence = queue.CurrentSequence + 1
		queue.stock <- qItem
		queue.CurrentSequence += 1
		queue.CurrentStock += 1
		queue.CurrentPid = qItem.Pid
		if overStock {
			log.Warn("[queue]", "over stock recover on ", qItem)
		}
	}
}

func (queue *Queue) Pop() *QueueItem {
	if queue != nil {
		qItem := <-queue.stock
		queue.CurrentStock -= 1
		log.Info("[queue]", "pop queue", qItem)
		return qItem
	}
	log.Error("[queue]", "error queue is nil")
	return nil
}
