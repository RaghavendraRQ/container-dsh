package pb

import (
	"context"

	v1 "github.com/raghavendrarq/container-dsh/api/gen/go/v1"
)

type ContainerService struct {
	v1.UnimplementedContainerServiceServer
}

func (c *ContainerService) GetContainerStats(ctx context.Context, req *v1.ContainerStatRequest) (*v1.Container, error) {
	return &v1.Container{}, nil

}
