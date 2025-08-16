package prom

import (
	"container-dsh/pkg/prom"
	"log"
)

func Run() {
	if err := prom.Run(); err != nil {
		log.Panic("Can't Run prometheus", err)
	}

}
