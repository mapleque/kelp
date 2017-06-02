package queue

import (
	"testing"
)

func TestQueueBase(t *testing.T) {
	queue := CreateQueue(10)
	qItem := queue.CreateItem("tmp_table", 0, "some interface")
	queue.Push(qItem)
	if queue.currentSequence != 1 {
		t.Error("error on currentSequence should be 1 but ", queue.currentSequence)
	}
	if queue.currentStock != 1 {
		t.Error("error on currentStock should be 1 but ", queue.currentStock)
	}

	oqItem := queue.Pop()
	if oqItem.sequence != 1 {
		t.Error("error on item sequence should be 1 but ", oqItem.sequence)
	}
	if oqItem.table != "tmp_table" || oqItem.pid != 0 || oqItem.data != "some interface" {
		t.Error("error on item data", oqItem.table, oqItem.pid, oqItem.data)
	}
}
