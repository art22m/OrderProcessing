package model

import "time"

type OrderState string

const (
	OrderStateCreated   = OrderState("Created")   // Заказ создан
	OrderStateProcessed = OrderState("Processed") // Заказ обработан
	OrderStateCompleted = OrderState("Completed") // Заказ выполнен
)

type Order struct {
	GoodsID         uint64
	WarehouseID     uint64
	DeliveryPointID uint64
	WorkerID        uint64
	Tracking        []OrderState
}

type OrderTracking struct {
	State OrderState
	Time  time.Time
}
