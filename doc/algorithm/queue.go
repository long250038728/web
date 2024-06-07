package algorithm

import "errors"

type queue interface {
	enqueue(value int32)
	dequeue() (int32, error)
}

//====================链式队列===========================

type queueLink struct {
	head *link
	curr *link
}

func (s *queueLink) enqueue(value int32) {
	l := &link{value: value}

	if s.head == nil {
		s.head = l
		s.curr = l
		return
	}
	l.next = s.head
	s.head = l
}

func (s *queueLink) dequeue() (int32, error) {
	if s.head == nil || s.curr == nil {
		return 0, errors.New("栈数据为空")
	}

	curr := s.head
	for {
		//只剩下一个元素
		if curr.next == nil {
			val := curr.value
			s.curr = nil
			s.head = nil
			return val, nil
		}

		//如果遍历当前项是tail时。就可以往前挪一个。
		//如果遍历当前项是nil时，需要往前再找一位，用到双向链表才能解决
		if curr.next == s.curr {
			val := s.curr.value
			curr.next = nil
			s.curr = curr
			return val, nil
		}
		curr = curr.next
	}
}

//====================数组队列===========================

type queueArray struct {
	head []int32
}

func (s *queueArray) enqueue(value int32) {
	if s.head == nil {
		s.head = make([]int32, 0, 10)
	}
	//自动扩容
	s.head = append([]int32{value}, s.head...)
}

func (s *queueArray) dequeue() (int32, error) {
	if s.head == nil || len(s.head) == 0 {
		return 0, errors.New("队列数据为空")
	}
	val := s.head[len(s.head)-1]
	s.head = s.head[:len(s.head)-1]
	return val, nil
}
