package cli

import "container-dsh/internal/container"

type Container struct {
	Id    uint
	Stats container.Stats
}

type Store struct {
	ContainerStats container.ContainersData
}
