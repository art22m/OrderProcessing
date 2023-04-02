package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"sync"

	"hw4/internal/model"
	completestep "hw4/internal/service/gateway/complete"
	createstep "hw4/internal/service/gateway/create"
	processstep "hw4/internal/service/gateway/process"
)

type Implementation struct {
	create   *createstep.Implementation
	process  *processstep.Implementation
	complete *completestep.Implementation
}

func New(create *createstep.Implementation, process *processstep.Implementation, complete *completestep.Implementation) *Implementation {
	return &Implementation{
		create:   create,
		process:  process,
		complete: complete,
	}
}

// Process реализует простую последовательную обработку заказа
func (i *Implementation) Process(workerID model.WorkerID, order model.Order) error {
	orderCreated, err := i.create.Create(workerID, order)
	if err != nil {
		return err
	}

	orderProcessed, err := i.process.Process(orderCreated)
	if err != nil {
		return err
	}

	orderCompleted, err := i.complete.Complete(orderProcessed)
	if err != nil {
		return err
	}

	orderInfo, err := prettyString(orderCompleted)
	if err != nil {
		return err
	}

	log.Println(orderInfo)

	return nil
}

// Pipeline реализует обработку заказов c приминением паттера Pipeline
func (i *Implementation) Pipeline(ctx context.Context, workerID model.WorkerID, orders <-chan model.Order) {
	createCh := i.create.Pipeline(ctx, workerID, orders)
	processCh := i.process.Pipeline(ctx, createCh)
	completeCh := i.complete.Pipeline(ctx, processCh)

	printCompletedOrders(completeCh)
}

// PipelineFan реализует обработку заказов c приминением паттернов Pipeline и FanIn-FanOut
func (i *Implementation) PipelineFan(ctx context.Context, workerID model.WorkerID, orders <-chan model.Order) {
	createCh := i.create.Pipeline(ctx, workerID, orders)

	const limit = 3
	fanOutProcess := make([]<-chan model.OrderPipeline, limit)
	for it := 0; it < limit; it++ {
		fanOutProcess[it] = i.process.Pipeline(ctx, createCh)
	}

	completeCh := i.complete.Pipeline(ctx, fanIn(ctx, fanOutProcess))

	printCompletedOrders(completeCh)
}

func fanIn(ctx context.Context, chans []<-chan model.OrderPipeline) <-chan model.OrderPipeline {
	multiplexed := make(chan model.OrderPipeline)

	var wg sync.WaitGroup
	for _, ch := range chans {
		wg.Add(1)

		go func(ch <-chan model.OrderPipeline) {
			defer wg.Done()
			for v := range ch {
				select {
				case <-ctx.Done():
					return

				case multiplexed <- v:
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(multiplexed)
	}()

	return multiplexed
}

// printCompletedOrders передает в стандартный поток вывода обработанные заказы в JSON формате структуры заказа
func printCompletedOrders(orders <-chan model.OrderPipeline) {
	for order := range orders {
		if order.Err != nil {
			log.Printf("Error while processing order with ID: [%d], err: [%v]", order.Order.ID, order.Err)
			continue
		}

		orderInfo, err := prettyString(order.Order)
		if err != nil {
			log.Printf("Error pretty printing order info: [%v]", err)
			continue
		}

		log.Println(orderInfo)
	}
}

// prettyString принимает model.Order и возвращает json в виде строки, построенный по заказу.
// Также добавляет табы в результирующую строку, чтобы сделать ее более читаемой.
func prettyString(order model.Order) (string, error) {
	str, err := json.Marshal(order)
	if err != nil {
		return "", nil
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}

	return prettyJSON.String(), nil
}
