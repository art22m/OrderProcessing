package create

import (
	"context"
	"time"

	"hw4/internal/model"
)

type Implementation struct{}

func New() *Implementation {
	return &Implementation{}
}

// Create присваивает айди воркера к заказу и добавляет в массив состояний новое состояние "Создан"
func (i *Implementation) Create(workerID model.WorkerID, order model.Order) (model.Order, error) {

	order.WorkerID = workerID
	order.Tracking = append(order.Tracking, model.OrderTracking{
		State: model.OrderStateCreated,
		Time:  time.Now().UTC(),
	})

	return order, nil
}

// Pipeline определяет начальный шаг всего пайплайна, принимает ID воркера и канал с заказами для обработки.
// Применяет к ним метод Create(...). Передает заказы в выходной канал, который является входным в следующий пайплайн.
func (i *Implementation) Pipeline(ctx context.Context, workerID model.WorkerID, orders <-chan model.Order) <-chan model.OrderPipeline {
	outCh := make(chan model.OrderPipeline)
	go func() {
		defer close(outCh)
		for order := range orders {
			orderR, err := i.Create(workerID, order)
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
