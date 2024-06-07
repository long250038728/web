package algorithm

import (
	"testing"
)

func Test_queueLink(t *testing.T) {
	stack := &queueLink{}
	stack.enqueue(1)
	stack.enqueue(2)
	stack.enqueue(3)
	stack.enqueue(4)
	stack.enqueue(5)

	for {
		v, err := stack.dequeue()
		if err != nil {
			return
		}
		t.Log(v)
	}
}

func Test_queueArray(t *testing.T) {
	stack := &queueArray{}
	stack.enqueue(1)
	stack.enqueue(2)
	stack.enqueue(3)
	stack.enqueue(4)
	stack.enqueue(5)

	for {
		v, err := stack.dequeue()
		if err != nil {
			return
		}
		t.Log(v)
	}
}
