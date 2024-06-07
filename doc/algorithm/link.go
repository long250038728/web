package algorithm

import "fmt"

type link struct {
	value int32
	next  *link
}

var head = &link{
	value: 1,
	next: &link{
		value: 2,
		next: &link{
			value: 3,
			next: &link{
				value: 4,
				next: &link{
					value: 5,
				},
			},
		},
	},
}

func reverse(head *link) *link {
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

func printLink(head *link) {
	if head == nil {
		return
	}
	fmt.Println(head.value)
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
