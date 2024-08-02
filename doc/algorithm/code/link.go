package code

type link struct {
	value int32
	next  *link
}

type doublyLink struct {
	value int32
	prev  *doublyLink
	next  *doublyLink
}

func reverse2(head *link) *link {
	// 有两个链 : 第一个是传入的链  第二个是空链
	// 把第一个链的第一个拿出来， 放到第二个链的第一个

	var prev, curr, next *link
	curr = head
	for curr != nil {
		next = curr.next
		curr.next = prev
		prev = curr
		curr = next
	}
	return prev
}

func reverse(head *link) *link {
	// 有两个链 : 第一个是传入的链  第二个新链
	// 把第一个链的第一个拿出来， 放到第二个链的第一个
	// 把第一个链的第一个拿出来， 放到第二个链的第一个
	var newLink *link

	for head != nil {
		//原链
		curr := head
		head = head.next

		//新链
		curr.next = newLink
		newLink = curr
	}
	return newLink
}

func hasCycle(head *link) bool {
	// 采用"龟兔赛跑"的思想

	if head == nil || head.next == nil {
		return false
	}

	one := head
	two := head.next
	for two != nil && two.next != nil {
		if one == two {
			return true
		}
		one = one.next
		two = two.next.next
	}
	return false
}

func printLink(head *link) {
	if head == nil {
		return
	}
	printLink(head.next)
}

//单链表反转
//链表中环的检测
//两个有序的链表合并
//删除链表倒数第n个结点
//求链表的中间结点

//利用哨兵简化实现难度
//重点留意边界条件处理
//举例画图，辅助思考
//多写多练，没有捷径

//链表中环的检测 （1,2,3,4,5,6,7,8,9,10）—— 使用步长不同比使用遍历对比首个性能会更好，同时特殊情况下会导致死循环
// 步长1 与 步长2:
// 这种情况需要等到10次
//	1 2 3 4 5       6 7 8 9 10
//	2 4 6 8 10      2 4 6 8 10
// 这种情况永远无法判断
//	1
//	2 4 6 8 10    2 4 6 8 10    2 4 6 8 10    2 4 6 8 10 ....
//
//  步长1 与 步长3:
//  这种情况需要等到5次
//	1 2 3 4 5
//  3 6 9 2 5
// 这种情况需要时7次(多2次)
//	1
//	3 6 9 2 5   8 1
//
//  步长1 与 步长4:
//  这种情况需要等到10次
//	1 2 3 4 5     6 7 8 9 10
//  4 8 2 6 10    4 8 2 6 10
// 这种情况永远无法判断
//	1
//	4 8 2 6 10    4 8 2 6 10    4 8 2 6 10    4 8 2 6 10  ....
