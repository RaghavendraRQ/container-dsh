package prom

import (
	"container-dsh/internal/container"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cpu_usage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cpu_usage",
			Help: "Shows cpu usage of the container",
		},
	)
)

func changeCpu() {
	cli := container.GetClient()
	containers, err := container.GetContainerList(cli)
	if err != nil {
		panic(err) // TODO: Want to propagate the error
	}
	targetContainer := containers[0] // For now, just take the first container

	for {

		targetContainerData, err := container.GetContainerData(cli, targetContainer)
		if err != nil {
			panic(err)
		}
		cpu_usage.Set(targetContainerData.CpuUsage)
		//fmt.Printf("CPU Usage for container %s: %f\n", targetContainer, targetContainerData.CpuUsage)
		time.Sleep(1 * time.Second)
	}

}

func Run() {
	prometheus.MustRegister(cpu_usage)
	http.Handle("/metrics", promhttp.Handler())
	go changeCpu()
	fmt.Println("Checking prometheus.(Server is running in the background)")
	log.Fatal(http.ListenAndServe(":3000", nil))

}
