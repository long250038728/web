package prometheus

import "testing"

func TestCount_do(t *testing.T) {
	count := NewCount()
	count.do()
	count.http()
}
