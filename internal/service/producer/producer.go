package producer

import (
	"math/rand"

	"hw4/internal/config"
	"hw4/internal/model"
)

// Orders возвращает канал, в котором будут лежать айдишники товаров для обработки.
// В текущей реализации айди товара определяется случайным образом. Количество задаётся
// значеним из конфига.
func Orders() <-chan model.GoodsID {
	result := make(chan model.GoodsID, config.OrdersNumber)
	go func() {
		defer close(result)

		for i := 0; i < config.OrdersNumber; i++ {
			result <- model.GoodsID(rand.Intn(100))
		}
	}()

	return result
}
