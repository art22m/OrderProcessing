package gateway

import (
	"hw4/internal/model"
	completestep "hw4/internal/service/gateway/complete"
	createstep "hw4/internal/service/gateway/create"
	processstep "hw4/internal/service/gateway/process"
)

type implementation struct {
	create   *createstep.Implementation
	process  *processstep.Implementation
	complete *completestep.Implementation
}

func New(create *createstep.Implementation, process *processstep.Implementation, complete *completestep.Implementation) *implementation {
	return &implementation{
		create:   create,
		process:  process,
		complete: complete,
	}
}

func (i *implementation) Process(goodsID model.GoodsID) error {
	orderCreated, err := i.create.Create(goodsID)
	if err != nil {
		return err
	}

	orderProcessed, err := i.process.Process(orderCreated)
	if err != nil {
		return err
	}

	if _, err = i.complete.Complete(orderProcessed); err != nil {
		return err
	}

	return nil
}
