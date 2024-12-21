package go_effective

import (
	"fmt"
	"strings"
	"testing"
)

func BenchmarkStr(b *testing.B) {
	str := ""
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			str = str + fmt.Sprintf("%d", j)
		}
	}
}
func BenchmarkBytes(b *testing.B) {
	builder := strings.Builder{}
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			builder.Write([]byte{byte(j)})
		}
	}
}

func BenchmarkGrowBytes(b *testing.B) {
	builder := strings.Builder{}
	builder.Grow(1000 * b.N)
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			builder.Write([]byte{byte(j)})
		}
	}
}

// 执行Benchmark进行性能分析，生成pprof文件 。通过 pprof工具分析pprof文件，生成火焰图及流程图进行分析是否有优化空间
// go test effective_benchamrk_test.go -bench=BenchmarkStr -benchmem -cpuprofile='cpu.pprof'
// go tool pprof -http :8889 cpu.pprof
