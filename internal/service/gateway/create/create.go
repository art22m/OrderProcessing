package create

import (
	"context"
	"hw4/internal/model"
	"time"
)

type Implementation struct {
	ids <-chan model.OrderID
}

func New(ids <-chan model.OrderID) *Implementation {
	return &Implementation{
		ids: ids,
	}
}

func (i *Implementation) Create(goodsID model.GoodsID) (model.Order, error) {
	orderID := <-i.ids
	order := model.Order{
		ID:      orderID,
		GoodsID: goodsID,
		Tracking: []model.OrderTracking{{
			State: model.OrderStateCreated,
			Time:  time.Now().UTC(),
		}},
	}

	return order, nil
}

func (i *Implementation) Pipeline(ctx context.Context, goodsIDCh <-chan model.GoodsID) <-chan model.PipelineOrder {
	outCh := make(chan model.PipelineOrder)
	go func() {
		defer close(outCh)
		for goodsID := range goodsIDCh {
			order, err := i.Create(goodsID)
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
