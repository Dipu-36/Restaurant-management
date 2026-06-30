package data

import (
	"Dipu-36/restaurant/internals/validator"
	"database/sql"
	"errors"
	"time"
)

type Restaurant struct {
	ID             int64     `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	StreetAddress  string    `json:"street_address"`
	OpeningTime    string    `json:"opening_time"`
	ClosingTime    string    `json:"closing_time"`
	DeliveryFee    int64     `json:"delivery_fee"`
	DeliveryRadius int32     `json:"delivery_radius"`
	IsOpen         bool      `json:"is_open"`
	Version        int32     `json:"version"`
}

func ValidateRestaurant(v *validator.Validator, restaurant *Restaurant) {

	v.Check(
		validator.NotBlank(restaurant.Name),
		"name",
		"must be provided",
	)

	v.Check(
		validator.MaxChars(restaurant.Name, 200),
		"name",
		"must not exceed 200 characters",
	)

	v.Check(
		validator.NotBlank(restaurant.Email),
		"email",
		"must be provided",
	)

	v.Check(
		validator.Matches(restaurant.Email, validator.EmailRX),
		"email",
		"must be a valid email address",
	)

	v.Check(
		validator.NotBlank(restaurant.Phone),
		"phone",
		"must be provided",
	)

	v.Check(
		validator.NotBlank(restaurant.StreetAddress),
		"street_address",
		"must be provided",
	)

	v.Check(
		restaurant.DeliveryFee >= 0,
		"delivery_fee",
		"must not be negative",
	)

	v.Check(
		restaurant.DeliveryRadius > 0,
		"delivery_radius",
		"must be greater than zero",
	)
	opening, err1 := time.Parse("15:04", restaurant.OpeningTime)
	closing, err2 := time.Parse("15:04", restaurant.ClosingTime)

	v.Check(
		err1 == nil,
		"opening_time",
		"must be in HH:MM format",
	)

	v.Check(
		err2 == nil,
		"closing_time",
		"must be in HH:MM format",
	)

	if err1 == nil && err2 == nil {
		v.Check(
			opening.Before(closing),
			"opening_time",
			"must be before closing time",
		)
	}
}

type RestaurantModel struct {
	DB *sql.DB
}

func (m RestaurantModel) Get() (*Restaurant, error) {

	query := `
		SELECT
			id,
			created_at,
			name,
			email,
			phone,
			street_address,
			opening_time,
			closing_time,
			delivery_fee,
			delivery_radius,
			is_open,
			version
		FROM restaurant
		LIMIT 1
	`

	var restaurant Restaurant

	err := m.DB.QueryRow(query).Scan(
		&restaurant.ID,
		&restaurant.CreatedAt,
		&restaurant.Name,
		&restaurant.Email,
		&restaurant.Phone,
		&restaurant.StreetAddress,
		&restaurant.OpeningTime,
		&restaurant.ClosingTime,
		&restaurant.DeliveryFee,
		&restaurant.DeliveryRadius,
		&restaurant.IsOpen,
		&restaurant.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &restaurant, nil
}

func (m RestaurantModel) Update(restaurant *Restaurant) error {

	query := `
		UPDATE restaurant
		SET
			name = $1,
			email = $2,
			phone = $3,
			street_address = $4,
			opening_time = $5,
			closing_time = $6,
			delivery_fee = $7,
			delivery_radius = $8,
			is_open = $9,
			version = version + 1
		WHERE id = $10
		AND version = $11
		RETURNING version
	`

	args := []interface{}{
		restaurant.Name,
		restaurant.Email,
		restaurant.Phone,
		restaurant.StreetAddress,
		restaurant.OpeningTime,
		restaurant.ClosingTime,
		restaurant.DeliveryFee,
		restaurant.DeliveryRadius,
		restaurant.IsOpen,
		restaurant.ID,
		restaurant.Version,
	}

	err := m.DB.QueryRow(query, args...).Scan(&restaurant.Version)

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
