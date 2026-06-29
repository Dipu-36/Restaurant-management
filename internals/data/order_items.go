package data

import (
	"Dipu-36/restaurant/internals/validator"
	"database/sql"
	"errors"
	"time"
)

type OrderItem struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`

	OrderID int64 `json:"order_id"`

	DishID int64 `json:"dish_id"`

	Quantity int32 `json:"quantity"`

	UnitPrice int64 `json:"unit_price"`

	Subtotal int64 `json:"subtotal"`

	Version int32 `json:"version"`
}

func ValidateOrderItem(v *validator.Validator, item *OrderItem) {

	v.Check(
		item.OrderID > 0,
		"order_id",
		"must be provided",
	)

	v.Check(
		item.DishID > 0,
		"dish_id",
		"must be provided",
	)

	v.Check(
		item.Quantity > 0,
		"quantity",
		"must be greater than zero",
	)

	v.Check(
		item.UnitPrice >= 0,
		"unit_price",
		"must not be negative",
	)

	v.Check(
		item.Subtotal >= 0,
		"subtotal",
		"must not be negative",
	)

	v.Check(
		item.Subtotal == int64(item.Quantity)*item.UnitPrice,
		"subtotal",
		"must equal quantity × unit price",
	)
}

type OrderItemModel struct {
	DB *sql.DB
}

func (m OrderItemModel) Insert(item *OrderItem) error {

	query := `
		INSERT INTO order_items
		(
			order_id,
			dish_id,
			quantity,
			unit_price,
			subtotal
		)
		VALUES
		(
			$1,$2,$3,$4,$5
		)
		RETURNING id, created_at, version
	`

	args := []interface{}{
		item.OrderID,
		item.DishID,
		item.Quantity,
		item.UnitPrice,
		item.Subtotal,
	}

	return m.DB.QueryRow(query, args...).Scan(
		&item.ID,
		&item.CreatedAt,
		&item.Version,
	)
}

func (m OrderItemModel) Get(id int64) (*OrderItem, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT
			id,
			created_at,
			order_id,
			dish_id,
			quantity,
			unit_price,
			subtotal,
			version
		FROM order_items
		WHERE id = $1
	`

	var item OrderItem

	err := m.DB.QueryRow(query, id).Scan(
		&item.ID,
		&item.CreatedAt,
		&item.OrderID,
		&item.DishID,
		&item.Quantity,
		&item.UnitPrice,
		&item.Subtotal,
		&item.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &item, nil
}

func (m OrderItemModel) GetByOrderID(orderID int64) ([]*OrderItem, error) {

	if orderID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT
			id,
			created_at,
			order_id,
			dish_id,
			quantity,
			unit_price,
			subtotal,
			version
		FROM order_items
		WHERE order_id = $1
		ORDER BY id
	`

	rows, err := m.DB.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*OrderItem{}

	for rows.Next() {

		var item OrderItem

		err := rows.Scan(
			&item.ID,
			&item.CreatedAt,
			&item.OrderID,
			&item.DishID,
			&item.Quantity,
			&item.UnitPrice,
			&item.Subtotal,
			&item.Version,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (m OrderItemModel) GetByDishID(dishID int64) ([]*OrderItem, error) {

	if dishID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT
			id,
			created_at,
			order_id,
			dish_id,
			quantity,
			unit_price,
			subtotal,
			version
		FROM order_items
		WHERE dish_id = $1
		ORDER BY created_at DESC
	`

	rows, err := m.DB.Query(query, dishID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*OrderItem{}

	for rows.Next() {

		var item OrderItem

		err := rows.Scan(
			&item.ID,
			&item.CreatedAt,
			&item.OrderID,
			&item.DishID,
			&item.Quantity,
			&item.UnitPrice,
			&item.Subtotal,
			&item.Version,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (m OrderItemModel) Update(item *OrderItem) error {

	query := `
		UPDATE order_items
		SET
			quantity = $1,
			unit_price = $2,
			subtotal = $3,
			version = version + 1
		WHERE id = $4
		AND version = $5
		RETURNING version
	`

	args := []interface{}{
		item.Quantity,
		item.UnitPrice,
		item.Subtotal,
		item.ID,
		item.Version,
	}

	err := m.DB.QueryRow(query, args...).Scan(&item.Version)

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

func (m OrderItemModel) Delete(id int64) error {

	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM order_items
		WHERE id = $1
	`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m OrderItemModel) GetAll() ([]*OrderItem, error) {

	query := `
		SELECT
			id,
			created_at,
			order_id,
			dish_id,
			quantity,
			unit_price,
			subtotal,
			version
		FROM order_items
		ORDER BY created_at DESC
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*OrderItem{}

	for rows.Next() {

		var item OrderItem

		err := rows.Scan(
			&item.ID,
			&item.CreatedAt,
			&item.OrderID,
			&item.DishID,
			&item.Quantity,
			&item.UnitPrice,
			&item.Subtotal,
			&item.Version,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
