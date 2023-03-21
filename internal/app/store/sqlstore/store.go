package sqlstore

import (
	"database/sql"

	"github.com/NoireHub/NATS-streaming-WebService/internal/app/store"
)

type Store struct {
	db              *sql.DB
	orderRepository *OrderRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Order() store.OrderRepository {
	if s.orderRepository != nil {
		return s.orderRepository
	}

	s.orderRepository = &OrderRepository{
		store: s,
	}

	return s.orderRepository
}
