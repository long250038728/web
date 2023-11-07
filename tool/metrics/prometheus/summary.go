package prometheus

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"time"
)

// Summary
type Summary struct {
	summary    prometheus.Summary
	summaryVec *prometheus.SummaryVec
}

func NewSummary() *Summary {
	opts := prometheus.SummaryOpts{
		Namespace:  "summary_ns",
		Subsystem:  "summary_ss",
		Name:       "summary_n",
		Help:       "this is my gauge",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}
	summary := prometheus.NewSummary(opts)

	opts2 := prometheus.SummaryOpts{
		Namespace:  "summary_ves_ns",
		Subsystem:  "summary_ves_ss",
		Name:       "summary_ves_n",
		Help:       "this is my summary_ves",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}
	summaryVec := prometheus.NewSummaryVec(opts2, []string{"method"})

	prometheus.MustRegister(summary, summaryVec)

	//count:       count_ns_count_ss_count_n 3
	//count_ves:   count_ves_ns_count_ves_ss_count_ves_n{one="one",two="two"} 50
	return &Summary{summary: summary, summaryVec: summaryVec}
}

func (g *Summary) do() {
	go func() {
		g.summary.Observe(0.05)
	}()

	go func() {
		for {
			start := time.Now()

			// 模拟请求处理时间为随机值
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

			t := time.Since(start).Seconds()

			fmt.Println(t)

			// 根据请求的方法，记录请求处理时间
			g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(t)

			time.Sleep(time.Second * 10)
		}
	}()
}

func (g *Summary) http() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8082", nil)
}
