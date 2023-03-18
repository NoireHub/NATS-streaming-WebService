package sqlstore

import "github.com/NoireHub/NATS-streaming-WebService/internal/app/model"

type OrderRepository struct {
}

func (r *OrderRepository) AddOrder(order *model.Order) error {
	return nil
}