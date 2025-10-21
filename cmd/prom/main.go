package prom

import (
	"context"
	"log"

	"github.com/raghavendrarq/container-dsh/pkg/prom"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := prom.Run(ctx); err != nil {
		log.Panic("Can't Run prometheus", err)
	}

}
