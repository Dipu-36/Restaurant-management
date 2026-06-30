package data

import (
	"Dipu-36/restaurant/internals/validator"
	"database/sql"
	"errors"
	"time"
)

const (
	PaymentProviderStripe = "stripe"
	PaymentProviderPayPal = "paypal"

	PaymentStatusPending   = "pending"
	PaymentStatusSucceeded = "succeeded"
	PaymentStatusFailed    = "failed"
	PaymentStatusRefunded  = "refunded"

	PaymentMethodCard      = "card"
	PaymentMethodApplePay  = "apple_pay"
	PaymentMethodGooglePay = "google_pay"
	PaymentMethodPayPal    = "paypal"
	PaymentMethodCash      = "cash"

	CurrencyUSD = "USD"
)

type Payment struct {
	ID            int64     `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	OrderID       int64     `json:"order_id"`
	Provider      string    `json:"provider"`
	TransactionID string    `json:"transaction_id"`
	PaymentStatus string    `json:"status"`
	Method        string    `json:"method"`
	Amount        int64     `json:"amount"`
	Currency      string    `json:"currency"`
	Version       int32     `json:"version"`
}

func ValidatePayment(v *validator.Validator, payment *Payment) {

	v.Check(
		payment.OrderID > 0,
		"order_id",
		"must be provided",
	)

	v.Check(
		validator.In(
			payment.Provider,
			PaymentProviderStripe,
			PaymentProviderPayPal,
		),
		"provider",
		"must contain a valid payment provider",
	)

	v.Check(
		validator.In(
			payment.PaymentStatus,
			PaymentStatusPending,
			PaymentStatusSucceeded,
			PaymentStatusFailed,
			PaymentStatusRefunded,
		),
		"status",
		"must contain a valid payment status",
	)

	v.Check(
		validator.In(
			payment.Method,
			PaymentMethodCard,
			PaymentMethodApplePay,
			PaymentMethodGooglePay,
			PaymentMethodPayPal,
			PaymentMethodCash,
		),
		"method",
		"must contain a valid payment method",
	)

	v.Check(
		payment.Amount > 0,
		"amount",
		"must be greater than zero",
	)

	v.Check(
		payment.Currency == CurrencyUSD,
		"currency",
		"must be USD",
	)
}

type PaymentModel struct {
	DB *sql.DB
}

func (m PaymentModel) Insert(payment *Payment) error {

	query := `
		INSERT INTO payments
		(
			order_id,
			provider,
			transaction_id,
			status,
			method,
			amount,
			currency
		)
		VALUES
		(
			$1,$2,$3,$4,$5,$6,$7
		)
		RETURNING id, created_at, version
	`

	args := []interface{}{
		payment.OrderID,
		payment.Provider,
		payment.TransactionID,
		payment.PaymentStatus,
		payment.Method,
		payment.Amount,
		payment.Currency,
	}

	return m.DB.QueryRow(query, args...).Scan(
		&payment.ID,
		&payment.CreatedAt,
		&payment.Version,
	)
}

func (m PaymentModel) Get(id int64) (*Payment, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT
			id,
			created_at,
			order_id,
			provider,
			transaction_id,
			status,
			method,
			amount,
			currency,
			version
		FROM payments
		WHERE id = $1
	`

	var payment Payment

	err := m.DB.QueryRow(query, id).Scan(
		&payment.ID,
		&payment.CreatedAt,
		&payment.OrderID,
		&payment.Provider,
		&payment.TransactionID,
		&payment.PaymentStatus,
		&payment.Method,
		&payment.Amount,
		&payment.Currency,
		&payment.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &payment, nil
}

func (m PaymentModel) GetByOrderID(orderID int64) (*Payment, error) {

	if orderID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT
			id,
			created_at,
			order_id,
			provider,
			transaction_id,
			status,
			method,
			amount,
			currency,
			version
		FROM payments
		WHERE order_id = $1
	`

	var payment Payment

	err := m.DB.QueryRow(query, orderID).Scan(
		&payment.ID,
		&payment.CreatedAt,
		&payment.OrderID,
		&payment.Provider,
		&payment.TransactionID,
		&payment.PaymentStatus,
		&payment.Method,
		&payment.Amount,
		&payment.Currency,
		&payment.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &payment, nil
}

func (m PaymentModel) UpdateStatus(payment *Payment) error {

	query := `
		UPDATE payments
		SET
			transaction_id = $1,
			status = $2,
			version = version + 1
		WHERE id = $3
		AND version = $4
		RETURNING version
	`

	args := []interface{}{
		payment.TransactionID,
		payment.PaymentStatus,
		payment.ID,
		payment.Version,
	}

	err := m.DB.QueryRow(query, args...).Scan(&payment.Version)

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
func (m PaymentModel) GetAll() ([]*Payment, error) {

	query := `
		SELECT
			id,
			created_at,
			order_id,
			provider,
			transaction_id,
			status,
			method,
			amount,
			currency,
			version
		FROM payments
		ORDER BY created_at DESC
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := []*Payment{}

	for rows.Next() {

		var payment Payment

		err := rows.Scan(
			&payment.ID,
			&payment.CreatedAt,
			&payment.OrderID,
			&payment.Provider,
			&payment.TransactionID,
			&payment.PaymentStatus,
			&payment.Method,
			&payment.Amount,
			&payment.Currency,
			&payment.Version,
		)
		if err != nil {
			return nil, err
		}

		payments = append(payments, &payment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (m PaymentModel) GetAllByStatus(status string) ([]*Payment, error) {

	query := `
		SELECT
			id,
			created_at,
			order_id,
			provider,
			transaction_id,
			status,
			method,
			amount,
			currency,
			version
		FROM payments
		WHERE ($1 = '' OR status = $1)
		ORDER BY created_at DESC
	`

	rows, err := m.DB.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := []*Payment{}

	for rows.Next() {

		var payment Payment

		err := rows.Scan(
			&payment.ID,
			&payment.CreatedAt,
			&payment.OrderID,
			&payment.Provider,
			&payment.TransactionID,
			&payment.PaymentStatus,
			&payment.Method,
			&payment.Amount,
			&payment.Currency,
			&payment.Version,
		)
		if err != nil {
			return nil, err
		}

		payments = append(payments, &payment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}
