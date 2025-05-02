package container

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type Container struct {
	ContainerID      string `json:"container_id"`
	ContainerImage   string `json:"container_image"`
	ContainerName    string `json:"container_name"`
	ContainerStatus  string `json:"container_status"`
	ContainerCreated string `json:"container_created"`
	ContainerNetwork string `json:"container_network"`
}

var (
	cli *client.Client
	ctx context.Context
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

func GetContainerData(cli *client.Client, containerId string) (*container.StatsResponseReader, error) {
	stats, err := cli.ContainerStatsOneShot(ctx, containerId)
	if err != nil {
		return nil, fmt.Errorf("error in getting container stats: %v", err)
	}
	return &stats, nil

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
