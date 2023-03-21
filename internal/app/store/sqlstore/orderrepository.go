package sqlstore

import "github.com/NoireHub/NATS-streaming-WebService/internal/app/model"

type OrderRepository struct {
	store *Store
	cache map[string]model.Order
}

func (r *OrderRepository) AddOrder(order *model.Order) error {
	return nil
}