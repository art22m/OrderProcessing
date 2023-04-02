package generator

import (
	"context"

	"hw4/internal/model"
)

// OrderIDs реализует паттерн генератор. Создаёт канал, в котором будут уникальные айдишники для заказов.
func OrderIDs(ctx context.Context) <-chan model.OrderID {
	result := make(chan model.OrderID)
	go func() {
		defer close(result)

		var counter uint64
		for {
			counter++

			select {
			case <-ctx.Done():
				return
			case result <- model.OrderID(counter):
			}
		}
	}()

	return result
}
