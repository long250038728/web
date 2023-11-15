package prometheus

import "testing"

func TestHistogram_do(t *testing.T) {
	histogram := NewHistogram()
	histogram.do()
	histogram.http()
}
