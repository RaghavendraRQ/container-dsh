package main

import (
	"flag"
	"fmt"
)

func main() {
	mode := flag.String("mode", "server", "Any one of the server, logger, cli, mock\nDefalut is Server")
	flag.Parse()

	switch *mode {
	case "server":
		fmt.Println("You have choosen server")
	case "logger":
		fmt.Println("You have choosen logger")
	case "cli":
		fmt.Println("You have choosen cli")
	case "mock":
		fmt.Println("You have choosen mock")
	default:
		fmt.Println("Can be ")
	}
}
