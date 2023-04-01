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

func (i *Implementation) Pipeline(ctx context.Context, workerID model.WorkerID, goodsIDCh <-chan model.GoodsID) <-chan model.PipelineOrder {
	outCh := make(chan model.PipelineOrder)
	go func() {
		defer close(outCh)
		for goodsID := range goodsIDCh {
			order, err := i.Create(workerID, goodsID)
			select {
			case <-ctx.Done():
				return

			case outCh <- model.PipelineOrder{
				Order:   order,
				GoodsID: goodsID,
				Err:     err,
			}:
			}
		}
	}()

	return outCh
}
