package data

import (
	"Dipu-36/restaurant/internals/validator"
	"database/sql"
	"errors"
	"time"
)

const (
	// Order OrderStatus
	OrderStatusPending   = "pending"
	OrderStatusPreparing = "preparing"
	OrderStatusReady     = "ready"
	OrderStatusCompleted = "completed"
	OrderStatusCancelled = "cancelled"

	// Order Type
	OrderTypeDelivery = "delivery"
	OrderTypePickup   = "pickup"
)

type Order struct {
	ID                int64     `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	CustomerID        int64     `json:"customer_id"`
	OrderStatus       string    `json:"status"`
	OrderType         string    `json:"order_type"`
	DeliveryAddressID int64     `json:"delivery_address_id,omitempty"`
	PickupTime        time.Time `json:"pickup_time,omitempty"`
	Subtotal          int64     `json:"subtotal"`
	DeliveryFee       int64     `json:"delivery_fee"`
	Tax               int64     `json:"tax"`
	Total             int64     `json:"total"`
	Version           int32     `json:"version"`
}

func ValidateOrder(v *validator.Validator, order *Order) {

	v.Check(
		order.CustomerID > 0,
		"customer_id",
		"must be provided",
	)

	v.Check(
		validator.In(
			order.OrderStatus,
			OrderStatusPending,
			OrderStatusPreparing,
			OrderStatusReady,
			OrderStatusCompleted,
			OrderStatusCancelled,
		),
		"status",
		"must contain a valid order status",
	)

	v.Check(
		validator.In(
			order.OrderType,
			OrderTypeDelivery,
			OrderTypePickup,
		),
		"order_type",
		"must be either delivery or pickup",
	)

	if order.OrderType == OrderTypeDelivery {
		v.Check(
			order.DeliveryAddressID > 0,
			"delivery_address_id",
			"must be provided for delivery orders",
		)

		v.Check(
			order.PickupTime.IsZero(),
			"pickup_time",
			"must not be provided for delivery orders",
		)
	}

	if order.OrderType == OrderTypePickup {
		v.Check(
			order.DeliveryAddressID == 0,
			"delivery_address_id",
			"must not be provided for pickup orders",
		)

		v.Check(
			!order.PickupTime.IsZero(),
			"pickup_time",
			"must be provided for pickup orders",
		)
	}

	v.Check(
		order.Subtotal >= 0,
		"subtotal",
		"must not be negative",
	)

	v.Check(
		order.DeliveryFee >= 0,
		"delivery_fee",
		"must not be negative",
	)

	v.Check(
		order.Tax >= 0,
		"tax",
		"must not be negative",
	)

	v.Check(
		order.Total >= 0,
		"total",
		"must not be negative",
	)

	v.Check(
		order.Total == order.Subtotal+order.DeliveryFee+order.Tax,
		"total",
		"must equal subtotal + delivery fee + tax",
	)
}

type OrderModel struct {
	DB *sql.DB
}

func (m OrderModel) Insert(order *Order) error {

	query := `
			INSERT INTO orders (
			customer_id,
			status,
			order_type,
			delivery_address_id,
			pickup_time,
			subtotal,
			delivery_fee,
			tax,
			total
		)
		VALUES (
			$1,$2,$3,$4,$5,
			$6,$7,$8,$9
		)
		RETURNING id, created_at, version
	`

	args := []interface{}{
		order.CustomerID,
		order.OrderStatus,
		order.OrderType,
		order.DeliveryAddressID,
		order.PickupTime,
		order.Subtotal,
		order.DeliveryFee,
		order.Tax,
		order.Total,
	}

	return m.DB.QueryRow(query, args...).Scan(
		&order.ID,
		&order.CreatedAt,
		&order.Version,
	)
}

func (m OrderModel) Get(id int64) (*Order, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT
			id,
			created_at,
			customer_id,
			status,
			order_type,
			delivery_address_id,
			pickup_time,
			subtotal,
			delivery_fee,
			tax,
			total,
			version
		FROM orders
		WHERE id = $1
	`

	var order Order

	err := m.DB.QueryRow(query, id).Scan(
		&order.ID,
		&order.CreatedAt,
		&order.CustomerID,
		&order.OrderStatus,
		&order.OrderType,
		&order.DeliveryAddressID,
		&order.PickupTime,
		&order.Subtotal,
		&order.DeliveryFee,
		&order.Tax,
		&order.Total,
		&order.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &order, nil
}

func (m OrderModel) UpdateStatus(order *Order) error {
	query := `
		UPDATE orders
		SET
			status = $1,
			version = version + 1
		WHERE id = $2
		AND version = $3
		RETURNING version
	`

	args := []interface{}{
		order.OrderStatus,
		order.ID,
		order.Version,
	}

	err := m.DB.QueryRow(query, args...).Scan(&order.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m OrderModel) GetCustomerOrders(customerID int64) ([]*Order, error) {

	query := `
		SELECT
			id,
			created_at,
			customer_id,
			status,
			order_type,
			delivery_address_id,
			pickup_time,
			subtotal,
			delivery_fee,
			tax,
			total,
			version
		FROM orders
		WHERE customer_id = $1
		ORDER BY created_at DESC
	`

	rows, err := m.DB.Query(query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []*Order{}

	for rows.Next() {
		var order Order

		err := rows.Scan(
			&order.ID,
			&order.CreatedAt,
			&order.CustomerID,
			&order.OrderStatus,
			&order.OrderType,
			&order.DeliveryAddressID,
			&order.PickupTime,
			&order.Subtotal,
			&order.DeliveryFee,
			&order.Tax,
			&order.Total,
			&order.Version,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (m OrderModel) GetAll(status string) ([]*Order, error) {

	query := `
		SELECT
			id,
			created_at,
			customer_id,
			status,
			order_type,
			delivery_address_id,
			pickup_time,
			subtotal,
			delivery_fee,
			tax,
			total,
			version
		FROM orders
		WHERE ($1 = '' OR status = $1)
		ORDER BY created_at DESC
	`

	rows, err := m.DB.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []*Order{}

	for rows.Next() {
		var order Order

		err := rows.Scan(
			&order.ID,
			&order.CreatedAt,
			&order.CustomerID,
			&order.OrderStatus,
			&order.OrderType,
			&order.DeliveryAddressID,
			&order.PickupTime,
			&order.Subtotal,
			&order.DeliveryFee,
			&order.Tax,
			&order.Total,
			&order.Version,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
