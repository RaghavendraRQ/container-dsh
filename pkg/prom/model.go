package prom

import (
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusMetrics struct {
	CpuUsage prometheus.Gauge
	MemUsage prometheus.Gauge
	NetIO    prometheus.Gauge
	DiskIO   prometheus.Gauge
}

func NewPrometheusMetrics() *PrometheusMetrics {
	return &PrometheusMetrics{
		CpuUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "container_dsh",
			Name:      "cpu_usage",
			Help:      "Shows CPU usage of the container",
		}),
		MemUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "container_dsh",
			Name:      "mem_usage",
			Help:      "Shows memory usage of the container",
		}),
		NetIO: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "container_dsh",
			Name:      "net_io",
			Help:      "Shows network I/O of the container",
		}),
		DiskIO: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "container_dsh",
			Name:      "disk_io",
			Help:      "Shows disk I/O of the container",
		}),
	}
}
