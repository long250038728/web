package prometheus

import (
	"testing"
)

func TestGauge_do(t *testing.T) {
	gauge := NewGauge()
	gauge.do()
	gauge.http()
}
