package code

import "testing"

func Test_arraySearch(t *testing.T) {
	index, err := binaryArraySearch([]int32{3, 4, 5, 6, 7}, 2)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(index)
}
