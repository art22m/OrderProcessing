package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

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

func (i *Implementation) Process(workerID model.WorkerID, goodsID model.GoodsID) error {
	orderCreated, err := i.create.Create(workerID, goodsID)
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

	fmt.Println(orderInfo)

	return nil
}

func (i *Implementation) Pipeline(ctx context.Context, workerID model.WorkerID, goodsIDCh <-chan model.GoodsID) {
	createCh := i.create.Pipeline(ctx, workerID, goodsIDCh)
	processCh := i.process.Pipeline(ctx, createCh)
	completeCh := i.complete.Pipeline(ctx, processCh)

	for order := range completeCh {
		if order.Err != nil {
			log.Printf("Error while processing order for goodsID: [%d], err: [%v]", order.GoodsID, order.Err)
			continue
		}

		orderInfo, err := prettyString(order.Order)
		if err != nil {
			log.Printf("Error pretty printing order info: [%v]", err)
			continue
		}

		fmt.Println(orderInfo)
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
