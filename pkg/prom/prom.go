package prom

import (
	"container-dsh/internal/container"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	testMetrics    = NewPrometheusMetrics()
	cli            = container.GetClient()
	scrapeInterval = 2 * time.Second // TODO: make configurable
)

func collectMetrics(containerId string) error {
	containerStats, err := container.GetContainerData(cli, containerId)
	if err != nil {
		return err
	}
	// container-dsh_METRICFAMILY{container_id, container_name} value
	testMetrics.CpuUsage.WithLabelValues(containerId, containerStats.Name).Set(containerStats.CpuUsage)
	testMetrics.MemUsage.WithLabelValues(containerId, containerStats.Name).Set(containerStats.MemUsage)
	testMetrics.DiskIO.WithLabelValues(containerId, containerStats.Name).Set(containerStats.DiskIO)
	testMetrics.NetIO.WithLabelValues(containerId, containerStats.Name).Set(containerStats.NetIO)

	return nil
}

func Run() error {
	prometheus.MustRegister(testMetrics.CpuUsage, testMetrics.DiskIO, testMetrics.NetIO, testMetrics.MemUsage)

	// Test for one sample container ID
	go func() {
		ticker := time.NewTicker(scrapeInterval)
		defer ticker.Stop()

		containerIds, err := container.GetContainerList(cli)
		if err != nil {
			log.Printf("Error getting container list: %v", err)
			return
		}

		for range ticker.C {
			for _, containerId := range containerIds {
				if err := collectMetrics(containerId); err != nil {
					log.Printf("Error collecting metrics for container %s: %v", containerId, err)
				}
			}
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Server is running on port 8080")
	return http.ListenAndServe(":8080", nil)
}
