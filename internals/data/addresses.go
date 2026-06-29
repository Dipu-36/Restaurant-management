package data

import (
	"Dipu-36/restaurant/internals/validator"
	"database/sql"
	"errors"
	"time"
)

type Address struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`

	CustomerID int64 `json:"customer_id"`

	StreetLine1 string `json:"street_line_1"`
	StreetLine2 string `json:"street_line_2,omitempty"`

	City string `json:"city"`

	State string `json:"state"`

	PostalCode string `json:"postal_code"`

	Country string `json:"country"`

	IsDefault bool `json:"is_default"`

	Version int32 `json:"version"`
}

func ValidateAddress(v *validator.Validator, address *Address) {

	v.Check(
		address.CustomerID > 0,
		"customer_id",
		"must be provided",
	)

	v.Check(
		validator.NotBlank(address.StreetLine1),
		"street_line_1",
		"must be provided",
	)

	v.Check(
		validator.MaxChars(address.StreetLine1, 200),
		"street_line_1",
		"must not exceed 200 characters",
	)

	v.Check(
		validator.MaxChars(address.StreetLine2, 200),
		"street_line_2",
		"must not exceed 200 characters",
	)

	v.Check(
		validator.NotBlank(address.City),
		"city",
		"must be provided",
	)

	v.Check(
		validator.MaxChars(address.City, 100),
		"city",
		"must not exceed 100 characters",
	)

	v.Check(
		validator.NotBlank(address.State),
		"state",
		"must be provided",
	)

	v.Check(
		validator.MaxChars(address.State, 100),
		"state",
		"must not exceed 100 characters",
	)

	v.Check(
		validator.NotBlank(address.PostalCode),
		"postal_code",
		"must be provided",
	)

	v.Check(
		validator.MaxChars(address.PostalCode, 20),
		"postal_code",
		"must not exceed 20 characters",
	)

	v.Check(
		validator.NotBlank(address.Country),
		"country",
		"must be provided",
	)

	v.Check(
		validator.MaxChars(address.Country, 100),
		"country",
		"must not exceed 100 characters",
	)
}

type AddressModel struct {
	DB *sql.DB
}

func (m AddressModel) Insert(address *Address) error {

	query := `
		INSERT INTO addresses
		(
			customer_id,
			street_line_1,
			street_line_2,
			city,
			state,
			postal_code,
			country,
			is_default
		)
		VALUES
		(
			$1,$2,$3,$4,
			$5,$6,$7,$8
		)
		RETURNING id, created_at, version
	`

	args := []interface{}{
		address.CustomerID,
		address.StreetLine1,
		address.StreetLine2,
		address.City,
		address.State,
		address.PostalCode,
		address.Country,
		address.IsDefault,
	}

	return m.DB.QueryRow(query, args...).Scan(
		&address.ID,
		&address.CreatedAt,
		&address.Version,
	)
}

func (m AddressModel) Update(address *Address) error {

	query := `
		UPDATE addresses
		SET
			street_line_1 = $1,
			street_line_2 = $2,
			city = $3,
			state = $4,
			postal_code = $5,
			country = $6,
			is_default = $7,
			version = version + 1
		WHERE id = $8
		AND version = $9
		RETURNING version
	`

	args := []interface{}{
		address.StreetLine1,
		address.StreetLine2,
		address.City,
		address.State,
		address.PostalCode,
		address.Country,
		address.IsDefault,
		address.ID,
		address.Version,
	}

	err := m.DB.QueryRow(query, args...).Scan(&address.Version)

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
