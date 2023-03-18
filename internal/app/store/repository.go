package store

import "github.com/NoireHub/NATS-streaming-WebService/internal/app/model"

type OrderRepository interface {
	AddOrder(*model.Order) error
}