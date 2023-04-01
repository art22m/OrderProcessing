package complete

import (
	"context"
	"hw4/internal/model"
	"time"
)

type implementation struct {
}

func New() *implementation {
	return &implementation{}
}

func (i *implementation) Complete(order model.Order) (model.Order, error) {
	// Процесс определения пвз
	time.Sleep(time.Second)

	order.DeliveryPointID = model.DeliveryPointID(uint64(order.GoodsID) + uint64(order.WarehouseID))
	order.Tracking = append(order.Tracking, model.OrderTracking{
		State: model.OrderStateCompleted,
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
