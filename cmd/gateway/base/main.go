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

	// Генератор айдишников заказа
	ids := generator.OrderIDs(ctx)

	// Операции над заказами
	create := createstep.New(ids)
	process := processstep.New()
	complete := completestep.New()

	// Канал с заказами для обработки
	orders := producer.Orders()

	// Сервер для обработки заказов
	server := gateway.New(create, process, complete)

	start := time.Now().UTC()

	// Это самая тривиальная обработка товаров, имеем одного рабочего
	var workerID model.WorkerID = 1
	for goodsID := range orders {
		if err := server.Process(workerID, goodsID); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Total duration %f seconds", time.Since(start).Seconds())
}
