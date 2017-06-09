package queue

import (
	"testing"
)

func TestQueueBase(t *testing.T) {
	queue := CreateQueue("tmp_task", 10)
	queue.Push("tmp_table", 0, "some interface")
	if queue.sequence != 1 {
		t.Error("error on current sequence should be 1 but ", queue.sequence)
	}
	if queue.stock != 1 {
		t.Error("error on current stock should be 1 but ", queue.stock)
	}

	oqItem := queue.Pop()
	if oqItem.Sequence != 1 {
		t.Error("error on item sequence should be 1 but ", oqItem.Sequence)
	}
	if oqItem.Name != "tmp_table" || oqItem.Flag != 0 || oqItem.Data != "some interface" {
		t.Error("error on item data", oqItem.Name, oqItem.Flag, oqItem.Data)
	}
}
