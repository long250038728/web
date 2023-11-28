package prometheus

import (
	"testing"
)

func TestNewSummary(t *testing.T) {
	summary := NewSummary()
	summary.do()
	summary.http()
}
