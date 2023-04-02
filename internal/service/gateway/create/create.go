package create

import (
	"context"
	"time"

	"hw4/internal/model"
)

type Implementation struct {
	ids <-chan model.OrderID
}

func New(ids <-chan model.OrderID) *Implementation {
	return &Implementation{
		ids: ids,
	}
}

// Create принимает айди воркера и айди товара, и по ним формирует новый заказ.
func (i *Implementation) Create(workerID model.WorkerID, goodsID model.GoodsID) (model.Order, error) {
	orderID := <-i.ids
	order := model.Order{
		ID:       orderID,
		GoodsID:  goodsID,
		WorkerID: workerID,
		Tracking: []model.OrderTracking{{
			State: model.OrderStateCreated,
			Time:  time.Now().UTC(),
		}},
	}

	return order, nil
}

// Pipeline определяет начальный шаг всего пайплайна, принимает айди воркера и канал с айдишниками товаров для обработки.
// По айдишнику, с помощью метода Create, формируется заказ и передается в канал, который является входным в следующий пайплайн.
func (i *Implementation) Pipeline(ctx context.Context, workerID model.WorkerID, goodsIDCh <-chan model.GoodsID) <-chan model.OrderPipeline {
	outCh := make(chan model.OrderPipeline)
	go func() {
		defer close(outCh)
		for goodsID := range goodsIDCh {
			order, err := i.Create(workerID, goodsID)
			select {
			case <-ctx.Done():
				return

			case outCh <- model.OrderPipeline{
				Order:   order,
				GoodsID: goodsID,
				Err:     err,
			}:
			}
		}
	}()

	return outCh
}
