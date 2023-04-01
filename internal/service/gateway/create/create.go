package create

import (
	"context"
	"hw4/internal/model"
	"time"
)

type implementation struct {
	ids <-chan model.OrderID
}

func New(ids <-chan model.OrderID) *implementation {
	return &implementation{
		ids: ids,
	}
}

func (i *implementation) Create(ctx context.Context, goodsID model.GoodsID) (model.Order, error) {
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

func (i *implementation) Pipeline(ctx context.Context, goodsIDCh <-chan model.GoodsID) <-chan model.PipelineOrder {
	outCh := make(chan model.PipelineOrder)
	go func() {
		defer close(outCh)
		for goodsID := range goodsIDCh {
			order, err := i.Create(ctx, goodsID)
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
