package main

import (
	"context"
	"log"
	"sync"
	"time"

	"hw4/internal/config"
	"hw4/internal/model"
	"hw4/internal/service/gateway"
	completestep "hw4/internal/service/gateway/complete"
	createstep "hw4/internal/service/gateway/create"
	processstep "hw4/internal/service/gateway/process"
	"hw4/internal/service/producer"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Операции над заказами
	create := createstep.New()
	process := processstep.New()
	complete := completestep.New()

	// Канал с заказами для обработки
	orders := producer.Orders()

	// Сервер для обработки заказов
	server := gateway.New(create, process, complete)

	start := time.Now().UTC()

	// Запуск воркеров
	var wg sync.WaitGroup
	for i := 0; i < config.WorkersNumber; i++ {
		wg.Add(1)
		go worker(ctx, model.WorkerID(i), &wg, server, orders)
	}

	wg.Wait()

	log.Printf("Total duration %f seconds", time.Since(start).Seconds())
}

func worker(ctx context.Context, workerID model.WorkerID, wg *sync.WaitGroup, server *gateway.Implementation, orders <-chan model.Order) {
	defer wg.Done()
	server.Pipeline(ctx, workerID, orders)
}
