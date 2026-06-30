package data

import (
	"Dipu-36/restaurant/internals/validator"
	"database/sql"
	"errors"
	"time"
)

type Category struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Name         string    `json:"name"`
	DisplayOrder int32     `json:"display_order"`
	Version      int32     `json:"version"`
}

func ValidateCategory(v *validator.Validator, category *Category) {
	v.Check(
		validator.NotBlank(category.Name),
		"name",
		"must be provided",
	)

	v.Check(
		validator.MaxChars(category.Name, 100),
		"name",
		"must not be more than 100 bytes long",
	)

	v.Check(
		category.DisplayOrder >= 0,
		"display_order",
		"must be greater than or equal to zero",
	)
}

type CategoryModel struct {
	DB DBTX
}

func (m CategoryModel) Insert(category *Category) error {
	query := `
		INSERT INTO categories (name, display_order)
		VALUES ($1, $2)
		RETURNING id, created_at, version
	`

	args := []interface{}{
		category.Name,
		category.DisplayOrder,
	}

	return m.DB.QueryRow(query, args...).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.Version,
	)
}

func (m CategoryModel) Get(id int64) (*Category, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, name, display_order, version
		FROM categories
		WHERE id = $1
	`

	var category Category

	err := m.DB.QueryRow(query, id).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.Name,
		&category.DisplayOrder,
		&category.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &category, nil
}

func (m CategoryModel) Update(category *Category) error {
	query := `
		UPDATE categories
		SET name = $1,
		    display_order = $2,
		    version = version + 1
		WHERE id = $3
		AND version = $4
		RETURNING version
	`

	args := []interface{}{
		category.Name,
		category.DisplayOrder,
		category.ID,
		category.Version,
	}

	err := m.DB.QueryRow(query, args...).Scan(&category.Version)

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

func (m CategoryModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM categories
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

func (m CategoryModel) GetAll() ([]*Category, error) {
	query := `
		SELECT id, created_at, name, display_order, version
		FROM categories
		ORDER BY display_order ASC, name ASC
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*Category{}

	for rows.Next() {
		var category Category

		err := rows.Scan(
			&category.ID,
			&category.CreatedAt,
			&category.Name,
			&category.DisplayOrder,
			&category.Version,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
