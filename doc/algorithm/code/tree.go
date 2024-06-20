package code

import "fmt"

type tree struct {
	val   int32
	left  *tree
	right *tree
}

// preTree 前序
func preTree(t *tree) {
	if t == nil {
		return
	}
	fmt.Println(t.val)
	preTree(t.left)
	preTree(t.right)
}

// inTree 中序
func inTree(t *tree) {
	if t == nil {
		return
	}
	inTree(t.left)
	fmt.Println(t.val)
	inTree(t.right)
}

// postTree 后序
func postTree(t *tree) {
	if t == nil {
		return
	}
	postTree(t.left)
	postTree(t.right)
	fmt.Println(t.val)
}

// BFS 广度优先遍历
func BFS(root *tree) {
	if root == nil {
		return
	}
	queue := []*tree{root} //入队

	// 先进先出的思想（队列）
	for len(queue) > 0 {
		node := queue[0] //出队
		queue = queue[1:]
		fmt.Println(node.val)

		if node.left != nil {
			queue = append(queue, node.left) //入队
		}
		if node.right != nil {
			queue = append(queue, node.right) //入队
		}
	}
}
