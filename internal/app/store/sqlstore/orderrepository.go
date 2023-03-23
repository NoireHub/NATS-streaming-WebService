package sqlstore

import (
	"fmt"

	"github.com/NoireHub/NATS-streaming-WebService/internal/app/model"
)

type OrderRepository struct {
	store *Store
	cache map[string]*model.Order
}

func (r *OrderRepository) AddOrder(order *model.Order) error {

	_, err := r.store.db.Exec(
		"INSERT INTO orders VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)",
		order.OrderUid,
		order.TrackNumber,
		order.Entry,
		order.Delivery,
		order.Payment,
		order.Items,
		order.Locale,
		order.InternalSignature,
		order.CustomerId,
		order.DeliveryService,
		order.ShardKey,
		order.SmId,
		order.DateCreated,
		order.OofShard,
	)
	if err != nil {
		return err
	}

	r.cache[order.OrderUid] = order

	return nil
}

func (r *OrderRepository) GetCache() error {
	rows, err := r.store.db.Queryx("SELECT * FROM orders")
	if err != nil {
		return err
	}
	
	for rows.Next() {
		var order model.Order
		if err:= rows.StructScan(&order); err != nil {
			return err
		}
		
		r.cache[order.OrderUid] = &order
	}

	return nil
}


func (r *OrderRepository) GetOrderById(orderUID string) (*model.Order, error) {
	if _, ok := r.cache[orderUID]; ok {
		return r.cache[orderUID], nil
	}

	order:= &model.Order{}
	if err:= r.store.db.QueryRowx(
		"SELECT * FROM orders WHERE order_uid = $1",
		orderUID,
	).StructScan(
		order,
	); err != nil {
		return nil, err
	}

	r.cache[order.OrderUid] = order

	fmt.Println(r.cache)

	return order, nil
}
