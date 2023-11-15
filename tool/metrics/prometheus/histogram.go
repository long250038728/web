package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// Histogram 累加值
type Histogram struct {
	histogram    prometheus.Histogram
	histogramVec *prometheus.HistogramVec
}

func NewHistogram() *Histogram {
	opts := prometheus.HistogramOpts{
		Namespace: "histogram_ns",
		Subsystem: "histogram_ss",
		Name:      "histogram_n",
		Help:      "this is my gauge",
		Buckets:   []float64{0.01, 0.03, 0.05, 0.1}, //需要递增
	}
	histogram := prometheus.NewHistogram(opts)

	opts2 := prometheus.HistogramOpts{
		Namespace: "histogram_ves_ns",
		Subsystem: "histogram_ves_ss",
		Name:      "histogram_ves_n",
		Help:      "this is my histogram_ves",
		Buckets:   []float64{0.01, 0.03, 0.05, 0.1}, //需要递增

		// 0 ~  0.01
		// 0 ~  0.03
		// 0 ~  0.05
		// 0 ~  0.1
	}
	histogramVec := prometheus.NewHistogramVec(opts2, []string{"one", "two"})

	prometheus.MustRegister(histogram, histogramVec)

	//histogram_ns_histogram_ss_histogram_n_bucket{le="0.01"} 1
	//histogram_ns_histogram_ss_histogram_n_bucket{le="0.03"} 2
	//histogram_ns_histogram_ss_histogram_n_bucket{le="0.05"} 3
	//histogram_ns_histogram_ss_histogram_n_bucket{le="0.1"} 4
	//histogram_ns_histogram_ss_histogram_n_bucket{le="+Inf"} 4
	//
	//histogram_ves_ns_histogram_ves_ss_histogram_ves_n_bucket{one="one",two="two",le="0.01"} 1
	//histogram_ves_ns_histogram_ves_ss_histogram_ves_n_bucket{one="one",two="two",le="0.03"} 2
	//histogram_ves_ns_histogram_ves_ss_histogram_ves_n_bucket{one="one",two="two",le="0.05"} 3
	//histogram_ves_ns_histogram_ves_ss_histogram_ves_n_bucket{one="one",two="two",le="0.1"} 4
	//histogram_ves_ns_histogram_ves_ss_histogram_ves_n_bucket{one="one",two="two",le="+Inf"} 4
	return &Histogram{histogram: histogram, histogramVec: histogramVec}
}

func (g *Histogram) do() {
	go func() {
		g.histogram.Observe(0.1)
		g.histogram.Observe(0.05)
		g.histogram.Observe(0.03)
		g.histogram.Observe(0.01)
	}()

	go func() {
		g.histogramVec.WithLabelValues("one", "two").Observe(0.1)
		g.histogramVec.WithLabelValues("one", "two").Observe(0.05)
		g.histogramVec.WithLabelValues("one", "two").Observe(0.03)
		g.histogramVec.WithLabelValues("one", "two").Observe(0.01)
	}()
}

func (g *Histogram) http() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8082", nil)
}
