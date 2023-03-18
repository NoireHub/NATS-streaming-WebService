package sqlstore

import (
	"database/sql"

	"github.com/NoireHub/NATS-streaming-WebService/internal/app/store"
)

type Store struct {
	db              *sql.DB
	orderRepository *OrderRepository
}

func (s *Store) Order() store.OrderRepository {
	return s.orderRepository
}
