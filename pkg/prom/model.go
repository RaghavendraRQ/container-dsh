package prom

import (
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusMetrics struct {
	CpuUsage *prometheus.GaugeVec
	MemUsage *prometheus.GaugeVec
	NetIO    *prometheus.GaugeVec
	DiskIO   *prometheus.GaugeVec
}

func NewPrometheusMetrics() *PrometheusMetrics {
	return &PrometheusMetrics{
		CpuUsage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "container_dsh",
			Name:      "cpu_usage",
			Help:      "Shows CPU usage of the container",
		}, []string{"container_id", "container_name"}),
		MemUsage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "container_dsh",
			Name:      "mem_usage",
			Help:      "Shows memory usage of the container",
		}, []string{"container_id", "container_name"}),
		NetIO: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "container_dsh",
			Name:      "net_io",
			Help:      "Shows network I/O of the container",
		}, []string{"container_id", "container_name"}),
		DiskIO: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "container_dsh",
			Name:      "disk_io",
			Help:      "Shows disk I/O of the container",
		}, []string{"container_id", "container_name"}),
	}
}
