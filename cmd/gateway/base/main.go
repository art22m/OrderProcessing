package main

import (
	"context"
	"fmt"
	"hw4/internal/service/gateway"
	completestep "hw4/internal/service/gateway/complete"
	createstep "hw4/internal/service/gateway/create"
	processstep "hw4/internal/service/gateway/process"
	"hw4/internal/service/generator"
	"hw4/internal/service/producer"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ids := generator.OrderIDs(ctx)

	create := createstep.New(ids)
	process := processstep.New()
	complete := completestep.New()

	server := gateway.New(create, process, complete)

	orders := producer.Orders()
	start := time.Now().UTC()
	for goodsID := range orders {
		if err := server.Process(goodsID); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Total duration: %f", time.Since(start).Seconds())
}
