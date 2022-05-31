package webserver

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

var (
	opsProcessed = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func NewMetricsRegistry() *prometheus.Registry {
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewGoCollector())

	//go counterInc()
	//reg.MustRegister(opsProcessed)

	return reg
}

func counterInc() {
	for {
		opsProcessed.Inc()
		time.Sleep(2 * time.Second)
	}
}
