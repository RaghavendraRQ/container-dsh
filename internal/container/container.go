package container

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

var (
	cli       *client.Client
	ctx       context.Context
	statsPool = sync.Pool{
		New: func() any {
			return &container.StatsResponse{}
		},
	}
)

func GetClient() *client.Client {
	if cli != nil {
		return cli
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	ctx = context.Background()
	return cli
}

func GetContainerList(cli *client.Client) ([]string, error) {
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("error in getting containers: %v", err)
	}
	var containerIds []string
	for _, cont := range containers {
		containerIds = append(containerIds, cont.ID)
	}
	return containerIds, nil
}

func GetContainerData(cli *client.Client, containerId string) (Container, error) {
	stats, err := cli.ContainerStatsOneShot(ctx, containerId)
	if err != nil {
		return Container{}, fmt.Errorf("error in getting container stats: %v", err)
	}
	defer stats.Body.Close()
	statsData := statsPool.Get().(*container.StatsResponse)
	if err := json.NewDecoder(stats.Body).Decode(statsData); err != nil {
		return Container{}, fmt.Errorf("error in decoding container stats: %v", err)
	}
	return Container{
		Stats:  NewMetrics(statsData),
		ID:     containerId,
		Name:   statsData.Name,
		Status: Running,
	}, nil

}

func GetImageList(cli *client.Client) ([]string, error) {
	images, err := cli.ImageList(ctx, image.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("error in getting images: %v", err)
	}
	var imageIds []string
	for _, img := range images {
		imageIds = append(imageIds, img.ID)
	}
	return imageIds, nil
}

func NewMetrics(statsData *container.StatsResponse) Stats {
	return Stats{
		CpuUsage: calculateCPUPercent(statsData),
		MemUsage: calculateMemUsage(statsData),
		NetIO:    calculateNetIO(statsData),
		DiskIO:   calculateDiskIO(statsData),
	}

}

func calculateCPUPercent(stat *container.StatsResponse) float64 {
	cpuDelta := float64(stat.CPUStats.CPUUsage.TotalUsage - stat.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stat.CPUStats.SystemUsage - stat.PreCPUStats.SystemUsage)
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		return (cpuDelta / systemDelta) * float64(len(stat.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return 0.0
}

func calculateMemUsage(stat *container.StatsResponse) float64 {
	return float64(stat.MemoryStats.Usage) / (1024 * 1024)
}

func calculateNetIO(stat *container.StatsResponse) float64 {
	if len(stat.Networks) > 0 {
		for _, net := range stat.Networks {
			return float64(net.RxBytes+net.TxBytes) / (1024 * 1024)
		}
	}
	return 0.0
}
func calculateDiskIO(stat *container.StatsResponse) float64 {
	if len(stat.BlkioStats.IoServiceBytesRecursive) > 0 {
		return float64(stat.BlkioStats.IoServiceBytesRecursive[0].Value) / (1024 * 1024)
	}
	return 0.0
}

func RunConainer(cli *client.Client, image string) {
	err := cli.ContainerStart(ctx, image, container.StartOptions{})
	if err != nil {
		fmt.Println("Can't able to start image")
	}
	log.Println("Started: ", image)

}
