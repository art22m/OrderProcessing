package complete

import (
	"context"
	"time"

	"hw4/internal/model"
)

type Implementation struct{}

func New() *Implementation {
	return &Implementation{}
}

// Complete вычисляет для сущности model.Order пункт выдачи заказов и добавляет финальный стейт в model.Order.Tracking
func (i *Implementation) Complete(order model.Order) (model.Order, error) {
	/* Здесь могла бы быть логика определения пункта выдачи заказов */
	time.Sleep(time.Second)

	order.DeliveryPointID = model.DeliveryPointID(uint64(order.GoodsID) + uint64(order.WarehouseID))
	order.Tracking = append(order.Tracking, model.OrderTracking{
		State: model.OrderStateCompleted,
		Time:  time.Now().UTC(),
	})

	return order, nil
}

// Pipeline определяет финальный шаг нашего пайплайна, принимает заказы со стейта <process> и применяет к ним метод
// Complete(order). Передает заказы в выходной канал, который является конечным каналом с обработанными заказами.
func (i *Implementation) Pipeline(ctx context.Context, orders <-chan model.OrderPipeline) <-chan model.OrderPipeline {
	outCh := make(chan model.OrderPipeline)
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

			case outCh <- model.OrderPipeline{
				Order: orderR,
				Err:   err,
			}:
			}
		}
	}()

	return outCh
}
