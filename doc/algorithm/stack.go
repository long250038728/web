package algorithm

import "errors"

type stack interface {
	enqueue(value int32)
	dequeue() (int32, error)
}

//====================链式栈===========================

type stackLink struct {
	head *link
}

func (s *stackLink) enqueue(value int32) {
	l := &link{value: value}

	if s.head == nil {
		s.head = l
		return
	}
	l.next = s.head
	s.head = l
}

func (s *stackLink) dequeue() (int32, error) {
	if s.head == nil {
		return 0, errors.New("栈数据为空")
	}
	val := s.head.value
	s.head = s.head.next
	return val, nil
}

//====================数组栈===========================

type stackArray struct {
	head []int32
}

func (s *stackArray) enqueue(value int32) {
	if s.head == nil {
		s.head = make([]int32, 0, 10)
	}
	//自动扩容
	s.head = append(s.head, value)
}

func (s *stackArray) dequeue() (int32, error) {
	if s.head == nil || len(s.head) == 0 {
		return 0, errors.New("栈数据为空")
	}
	val := s.head[len(s.head)-1]
	s.head = s.head[:len(s.head)-1]
	return val, nil
}
