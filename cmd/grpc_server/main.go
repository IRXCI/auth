package main

import (
	"context"
	"log"

	"github.com/IRXCI/auth/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err)
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err)
	}
}
