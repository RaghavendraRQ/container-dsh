package pb

import "github.com/raghavendrarq/container-dsh/pkg/pb"

func Run() {
	if err := pb.Run(); err != nil {
		panic(err)
	}
}
