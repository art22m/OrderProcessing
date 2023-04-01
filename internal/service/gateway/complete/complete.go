package complete

import (
	"context"
	"time"

	"hw4/internal/model"
)

type Implementation struct {
}

func New() *Implementation {
	return &Implementation{}
}

func (i *Implementation) Complete(order model.Order) (model.Order, error) {
	// Процесс определения пвз
	time.Sleep(time.Second)

	order.DeliveryPointID = model.DeliveryPointID(uint64(order.GoodsID) + uint64(order.WarehouseID))
	order.Tracking = append(order.Tracking, model.OrderTracking{
		State: model.OrderStateCompleted,
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

			orderR, err := i.Complete(order.Order)
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
