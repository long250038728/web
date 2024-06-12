package algorithm

import (
	"fmt"
	"testing"
)

func Test_Reverse(t *testing.T) {
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
	printLink(reverse(head))
}

func Test_hasCycle(t *testing.T) {
	var head = &link{value: 1, next: nil}

	head.next = &link{
		value: 2,
		next: &link{
			value: 3,
			next: &link{
				value: 4,
				next: &link{
					value: 5,
					next:  head,
				},
			},
		},
	}

	fmt.Println(hasCycle(head))
}
