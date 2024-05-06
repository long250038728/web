package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// Gauge 当前值
type Gauge struct {
	gauge    prometheus.Gauge
	gaugeVec *prometheus.GaugeVec
}

func NewGauge() *Gauge {
	//prometheus.NewGauge()
	//prometheus.NewGaugeVec()

	//prometheus.NewCounter()
	//prometheus.NewCounterVec()

	//prometheus.NewSummary()
	//prometheus.NewSummaryVec()

	opts := prometheus.GaugeOpts{
		Namespace: "gauge_ns",
		Subsystem: "gauge_ss",
		Name:      "gauge_n",
		Help:      "this is my gauge",
	}
	gauge := prometheus.NewGauge(opts)

	opts2 := prometheus.GaugeOpts{
		Namespace: "gauge_ves_ns",
		Subsystem: "gauge_ves_ss",
		Name:      "gauge_ves_n",
		Help:      "this is my gauge_ves",
	}
	gaugeVec := prometheus.NewGaugeVec(opts2, []string{"first", "second"})

	prometheus.MustRegister(gauge, gaugeVec)

	//gauge:       gauge_ns_gauge_ss_gauge_n 56
	//gauge_ves:   gauge_ves_ns_gauge_ves_ss_gauge_ves_n{first="first",second="second"} 50
	return &Gauge{gauge: gauge, gaugeVec: gaugeVec}
}

func (g *Gauge) do() {
	go func() {
		g.gauge.Set(50)
		g.gauge.Add(10)
		g.gauge.Sub(5)
		g.gauge.Dec()
		g.gauge.Inc()
		g.gauge.Inc()
	}()

	go func() {
		g.gaugeVec.WithLabelValues("first", "second").Add(10)
		g.gaugeVec.WithLabelValues("first", "second").Add(20)

		g.gaugeVec.WithLabelValues("first", "second").Add(10)
		g.gaugeVec.WithLabelValues("first", "second").Add(10)
	}()
}

func (g *Gauge) http() {
	http.Handle("/metrics", promhttp.Handler())
	_ = http.ListenAndServe(":8081", nil)
}
