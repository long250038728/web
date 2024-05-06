package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

// Count 累加值
type Count struct {
	count    prometheus.Counter
	countVec *prometheus.CounterVec
}

func NewCount() *Count {
	opts := prometheus.CounterOpts{
		Namespace: "count_ns",
		Subsystem: "count_ss",
		Name:      "count_n",
		Help:      "this is my gauge",
	}
	count := prometheus.NewCounter(opts)

	opts2 := prometheus.CounterOpts{
		Namespace: "count_ves_ns",
		Subsystem: "count_ves_ss",
		Name:      "count_ves_n",
		Help:      "this is my count_ves",
	}
	countVec := prometheus.NewCounterVec(opts2, []string{"one", "two"})

	prometheus.MustRegister(count, countVec)

	//count:       count_ns_count_ss_count_n 3
	//count_ves:   count_ves_ns_count_ves_ss_count_ves_n{one="one",two="two"} 50
	return &Count{count: count, countVec: countVec}
}

func (g *Count) do() {
	go func() {
		for {
			g.count.Inc()
			g.count.Desc()
			g.count.Inc()
			g.count.Inc()
			time.Sleep(time.Second)
		}
	}()

	go func() {
		g.countVec.WithLabelValues("one", "two").Add(10)
		g.countVec.WithLabelValues("one", "two").Add(20)

		g.countVec.WithLabelValues("one", "two").Add(10)
		g.countVec.WithLabelValues("one", "two").Add(10)
	}()
}

func (g *Count) http() {
	http.Handle("/metrics", promhttp.Handler())
	_ = http.ListenAndServe(":8082", nil)
}
