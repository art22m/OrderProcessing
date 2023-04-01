package generator

import (
	"context"

	"hw4/internal/model"
)

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
