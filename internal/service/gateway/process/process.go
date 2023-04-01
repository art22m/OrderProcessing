package process

import (
	"context"
	"hw4/internal/model"
	"time"
)

type Implementation struct {
}

func New() *Implementation {
	return &Implementation{}
}

func (i *Implementation) Process(order model.Order) (model.Order, error) {
	// Обработка
	time.Sleep(3 * time.Second)

	order.WarehouseID = model.WarehouseID(order.GoodsID % 2)
	order.Tracking = append(order.Tracking, model.OrderTracking{
		State: model.OrderStateProcessed,
		Time:  time.Now().UTC(),
	})

	return order, nil
}

func (i *Implementation) Pipeline(ctx context.Context, orders <-chan model.PipelineOrder) <-chan model.PipelineOrder {
	outCh := make(chan model.PipelineOrder)
	go func() {
		defer close(outCh)
		for order := range orders {
			if order.Err != nil {
				select {
				case <-ctx.Done():
					return

				case outCh <- order:
				}
			}

			orderR, err := i.Process(order.Order)
			select {
			case <-ctx.Done():
				return

			case outCh <- model.PipelineOrder{
				Order: orderR,
				Err:   err,
			}:
			}
		}
	}()

	return outCh
}
