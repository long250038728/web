package prometheus

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

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
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}, //绝对值偏差
	}
	summaryVec := prometheus.NewSummaryVec(opts2, []string{"method"})

	//summary_ns_summary_ss_summary_n{quantile="0.5"} 0.05
	//summary_ns_summary_ss_summary_n{quantile="0.9"} 0.05
	//summary_ns_summary_ss_summary_n{quantile="0.99"} 0.05
	//summary_ves_ns_summary_ves_ss_summary_ves_n{method="GET",quantile="0.5"} 0.05
	//summary_ves_ns_summary_ves_ss_summary_ves_n{method="GET",quantile="0.9"} 0.05
	//summary_ves_ns_summary_ves_ss_summary_ves_n{method="GET",quantile="0.99"} 0.05
	prometheus.MustRegister(summary, summaryVec)

	return &Summary{summary: summary, summaryVec: summaryVec}
}

func (g *Summary) do() {
	go func() {
		//summary_ns_summary_ss_summary_n{quantile="0.5"} 106  //平均数
		//summary_ns_summary_ss_summary_n{quantile="0.9"} 110
		//summary_ns_summary_ss_summary_n{quantile="0.99"} 140
		//summary_ns_summary_ss_summary_n_sum 1183
		//summary_ns_summary_ss_summary_n_count 11

		g.summary.Observe(110)
		g.summary.Observe(100)
		g.summary.Observe(98)
		g.summary.Observe(100)
		g.summary.Observe(105)
		g.summary.Observe(108)
		g.summary.Observe(106)
		g.summary.Observe(106)
		g.summary.Observe(100)
		g.summary.Observe(110)
		g.summary.Observe(140)

		// 98  100 100 100 105 106 106 108 110  110 140
	}()

	go func() {
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(0.01)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(0.02)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(0.03)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(0.04)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(0.05)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(0.06)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(0.07)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(0.08)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(0.09)

		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(0.11)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(0.12)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(0.13)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(0.14)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(0.15)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(0.16)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(0.17)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(0.18)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(0.19)

		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(1.11)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(2.12)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(3.13)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(4.14)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(5.15)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(6.16)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(7.17)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(8.18)
		g.summaryVec.With(prometheus.Labels{"method": "GET"}).Observe(9.19)

		//{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}
		//summary_ves_ns_summary_ves_ss_summary_ves_n{method="GET",quantile="0.5"} 0.15   中位数0.15
		//summary_ves_ns_summary_ves_ss_summary_ves_n{method="GET",quantile="0.9"} 7.17   九分位数 7.17
		//summary_ves_ns_summary_ves_ss_summary_ves_n{method="GET",quantile="0.99"} 9.19  九九分位  9.19
	}()
}

func (g *Summary) http() {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":8083", nil)
	fmt.Println(err)
}
