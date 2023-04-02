package process

import (
	"context"
	"time"

	"hw4/internal/model"
)

type Implementation struct{}

func New() *Implementation {
	return &Implementation{}
}

// Process вычисляет ID склада для заказа и добавляет в массив состояний новое состояние "Обработан"
func (i *Implementation) Process(order model.Order) (model.Order, error) {
	/* Здесь могла бы быть логика определения склада */
	time.Sleep(2 * time.Second)

	order.WarehouseID = model.WarehouseID(order.GoodsID % 2)
	order.Tracking = append(order.Tracking, model.OrderTracking{
		State: model.OrderStateProcessed,
		Time:  time.Now().UTC(),
	})

	return order, nil
}

// Pipeline определяет промежуточный шаг нашего пайплайна, принимает заказы со стейта <create> и применяет к ним метод
// Process(order). Передает заказы в выходной канал, который является входным в следующий пайплайн.
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

			orderR, err := i.Process(order.Order)
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
