package main

import (
	"context"
	"fmt"
	"hw4/internal/model"
	"log"
	"time"

	"hw4/internal/service/gateway"
	completestep "hw4/internal/service/gateway/complete"
	createstep "hw4/internal/service/gateway/create"
	processstep "hw4/internal/service/gateway/process"
	"hw4/internal/service/generator"
	"hw4/internal/service/producer"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Id generator
	ids := generator.OrderIDs(ctx)

	// Operations
	create := createstep.New(ids)
	process := processstep.New()
	complete := completestep.New()

	// Order to process
	orders := producer.Orders()

	// Server
	server := gateway.New(create, process, complete)

	start := time.Now().UTC()

	var workerID model.WorkerID = 1
	for goodsID := range orders {
		if err := server.Process(workerID, goodsID); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Total duration: %f", time.Since(start).Seconds())
}
