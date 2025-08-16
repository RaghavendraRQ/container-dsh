package prom

import (
	"container-dsh/internal/container"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	testMetrics = NewPrometheusMetrics()
	cli         = container.GetClient()
)

func collectMetrics(containerId string) error {
	containerStats, err := container.GetContainerData(cli, containerId)
	if err != nil {
		return err
	}
	// container_stats{container_id, container_name} value
	testMetrics.CpuUsage.WithLabelValues(containerId, containerStats.Name).Set(containerStats.CpuUsage)
	testMetrics.MemUsage.WithLabelValues(containerId, containerStats.Name).Set(rand.Float64())
	testMetrics.DiskIO.WithLabelValues(containerId, containerStats.Name).Set(containerStats.DiskIO)
	testMetrics.NetIO.WithLabelValues(containerId, containerStats.Name).Set(containerStats.NetIO)

	return nil
}

func Run() error {
	prometheus.MustRegister(testMetrics.CpuUsage, testMetrics.DiskIO, testMetrics.NetIO, testMetrics.MemUsage)

	// Test for one sample container ID
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		containerId := "3cd562e2b938"

		for range ticker.C {
			if err := collectMetrics(containerId); err != nil {
				log.Printf("Error collecting metrics for container %s: %v", containerId, err)
			}
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Server is running on port 8080")
	return http.ListenAndServe(":8080", nil)
}
