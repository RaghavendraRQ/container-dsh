package main

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
    ctx := context.Background()
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        log.Fatalln("Error: ", err)
    }

    defer cli.Close()

    containers, err := cli.ContainerList(ctx, container.ListOptions{All:true})
    if err != nil {
        log.Fatalln("Error: ", err)
    }

    for _, container := range containers {
        fmt.Printf("%s %s (status: %s)\n, %s", container.ID, container.Image, container.Status, container.State)
    }




}
