package main

import (
	"container-dsh/internal/container"
	"fmt"
)

func main() {
	cli := container.GetClient()
	status, _ := container.GetStatusById(cli, "69e139e71e16")

	fmt.Println("Container Status: ", status)

}
