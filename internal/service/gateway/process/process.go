package process

import (
	"context"
	"hw4/internal/model"
	"time"
)

type implementation struct {
	workerID model.WorkerID
}

func New(workerID model.WorkerID) *implementation {
	return &implementation{
		workerID: workerID,
	}
}

func (i *implementation) Process(ctx context.Context, order model.Order) (model.Order, error) {
	// Воркер берет заказ для обработки
	order.WorkerID = i.workerID

	// Обрабатывает
	time.Sleep(3 * time.Second)

	// Отсортировал на склад
	order.WarehouseID = model.WarehouseID(order.GoodsID % 2)
	order.Tracking = append(order.Tracking, model.OrderTracking{
		State: model.OrderStateProcessed,
		Time:  time.Now().UTC(),
	})

	return order, nil
}

func (i *implementation) Pipeline(ctx context.Context, orders <-chan model.PipelineOrder) <-chan model.PipelineOrder {
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

			orderR, err := i.Process(ctx, order.Order)
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
