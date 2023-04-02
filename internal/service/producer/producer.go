package producer

import (
	"hw4/internal/config"
	"hw4/internal/model"
	"math/rand"
)

// Orders возвращает канал, содержащий структуры заказов с инициализированным полем ID товара и ID заказа.
// В текущей реализации ID товара определяется случайным образом.
// Количество задаётся значеним из конфига.
func Orders() <-chan model.Order {
	result := make(chan model.Order, config.OrdersNumber)
	go func() {
		defer close(result)

		for i := 0; i < config.OrdersNumber; i++ {
			result <- model.Order{
				ID:      model.OrderID(i),
				GoodsID: model.GoodsID(rand.Intn(100)),
			}
		}
	}()

	return result
}
