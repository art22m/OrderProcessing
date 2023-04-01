package producer

import (
	"hw4/internal/config"
	"hw4/internal/model"
	"math/rand"
)

func Orders() <-chan model.GoodsID {
	result := make(chan model.GoodsID, config.GoodsNumber)
	go func() {
		defer close(result)

		for i := 0; i < config.GoodsNumber; i++ {
			result <- model.GoodsID(rand.Int())
		}
	}()

	return result
}
