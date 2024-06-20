package code

import "testing"

var tr = &tree{
	val: 1,
	left: &tree{
		val: 2,
		left: &tree{
			val: 4,
		},
		right: &tree{
			val: 5,
			left: &tree{
				val: 8,
			},
		},
	},
	right: &tree{
		val: 3,
		left: &tree{
			val: 6,
		},
		right: &tree{
			val: 7,
		},
	},
}

//   					  1
//			 2    						3
//		 4         5			   6	     7
//		       8

func Test_preTree(t *testing.T) {
	// 1 2 4 5 8  3 6 7
	preTree(tr) // 先打印这个节点，然后再打印它的左子树，最后打印它的右子树。
}

func Test_inTree(t *testing.T) {
	// 4 2 8 5 1 6 3 7
	inTree(tr) // 先打印它的左子树，然后再打印它本身，最后打印它的右子树。
}

func Test_postTree(t *testing.T) {
	// 4 8 5 2 6 7 3 1
	postTree(tr) // 先打印它的左子树，然后再打印它的右子树，最后打印这个节点本身。
}

func Test_BFS(t *testing.T) {
	BFS(tr)
}
