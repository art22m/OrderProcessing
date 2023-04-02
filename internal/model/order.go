package model

import "time"

type OrderState string

const (
	OrderStateCreated   = OrderState("Created")   // Заказ создан
	OrderStateProcessed = OrderState("Processed") // Заказ обработан
	OrderStateCompleted = OrderState("Completed") // Заказ выполнен
)

type OrderID uint64

type Order struct {
	ID              OrderID
	GoodsID         GoodsID
	WarehouseID     WarehouseID
	DeliveryPointID DeliveryPointID
	WorkerID        WorkerID
	Tracking        []OrderTracking
}

type OrderTracking struct {
	State OrderState
	Time  time.Time
}

type OrderPipeline struct {
	Order   Order
	GoodsID GoodsID
	Err     error
}
