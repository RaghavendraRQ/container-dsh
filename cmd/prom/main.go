package prom

import (
	"container-dsh/pkg/prom"
	"context"
	"log"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := prom.Run(ctx); err != nil {
		log.Panic("Can't Run prometheus", err)
	}

}
