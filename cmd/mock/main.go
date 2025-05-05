package main

import (
	"container-dsh/internal/container"
	"fmt"
)

func main() {
	cli := container.GetClient()
	containers, err := container.GetContainerList(cli)
	if err != nil {
		panic(err)
	}
	for _, containerId := range containers {
		fmt.Println("Container ID: ", containerId)
	}
	images, _ := container.GetImageList(cli)

	for _, image := range images {
		fmt.Println("image ID: ", image)
	}

}
