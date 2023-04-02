package main

import (
	"log"
	"time"

	"hw4/internal/model"
	"hw4/internal/service/gateway"
	completestep "hw4/internal/service/gateway/complete"
	createstep "hw4/internal/service/gateway/create"
	processstep "hw4/internal/service/gateway/process"
	"hw4/internal/service/producer"
)

func main() {
	// Операции над заказами
	create := createstep.New()
	process := processstep.New()
	complete := completestep.New()

	// Канал с заказами для обработки
	orders := producer.Orders()

	// Сервер для обработки заказов
	server := gateway.New(create, process, complete)

	start := time.Now().UTC()

	// Это самая тривиальная обработка товаров, имеем одного рабочего
	var workerID model.WorkerID = 0
	for order := range orders {
		if err := server.Process(workerID, order); err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Total duration %f seconds", time.Since(start).Seconds())
}
