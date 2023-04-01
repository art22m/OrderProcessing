package main

import (
	"context"
	"fmt"
	"hw4/internal/config"
	"hw4/internal/model"
	"sync"
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

	// Orders to process
	orders := producer.Orders()

	// Server
	server := gateway.New(create, process, complete)

	start := time.Now().UTC()

	var wg sync.WaitGroup
	for i := 0; i < config.WorkersNumber; i++ {
		wg.Add(1)
		go worker(ctx, model.WorkerID(i), &wg, server, orders)
	}

	wg.Wait()

	fmt.Printf("Total duration: %f", time.Since(start).Seconds())
}

func worker(ctx context.Context, workerID model.WorkerID, wg *sync.WaitGroup, server *gateway.Implementation, orders <-chan model.GoodsID) {
	defer wg.Done()
	server.Pipeline(ctx, workerID, orders)
}
