package algorithm

import (
	"testing"
)

func Test_stackLink(t *testing.T) {
	stack := &stackLink{}
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

func Test_stackArray(t *testing.T) {
	stack := &stackArray{}
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
