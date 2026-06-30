package data

import (
	"Dipu-36/restaurant/internals/validator"
	"database/sql"
	"errors"
	"time"
)

// TO DO : THink about the data or the dish like what to keep and what not to keep
type Dish struct {
	ID              int64     `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Price           int64     `json:"price"`
	CategoryID      int64     `json:"category_id"`
	ImageURL        string    `json:"image_url,omitempty"`
	IsAvailable     bool      `json:"is_available"`
	IsVegetarian    bool      `json:"is_vegetarian"`
	IsFeatured      bool      `json:"is_featured"`
	PreparationTime int32     `json:"preparation_time"`
	Version         int32     `json:"version"`
}

func ValidateDish(v *validator.Validator, dish *Dish) {
	v.Check(
		validator.NotBlank(dish.Name),
		"name",
		"must be provided",
	)

	v.Check(
		validator.MaxChars(dish.Name, 200),
		"name",
		"must not be more than 200 bytes long",
	)

	v.Check(
		validator.NotBlank(dish.Description),
		"description",
		"must be provided",
	)

	v.Check(
		validator.MaxChars(dish.Description, 1000),
		"description",
		"must not be more than 1000 bytes long",
	)

	v.Check(
		dish.Price > 0,
		"price",
		"must be greater than zero",
	)

	v.Check(
		dish.CategoryID > 0,
		"category_id",
		"must be provided",
	)

	v.Check(
		validator.Between(int(dish.PreparationTime), 1, 180),
		"preparation_time",
		"must be between 1 and 180 minutes",
	)
}

type DishModel struct {
	DB DBTX
}

func (m DishModel) Insert(dish *Dish) error {

	query := `
		INSERT INTO dishes
		(name, description, price, category_id, image_url,
		is_available, is_vegetarian, is_featured, preparation_time)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, created_at, version
	`

	args := []interface{}{
		dish.Name,
		dish.Description,
		dish.Price,
		dish.CategoryID,
		dish.ImageURL,
		dish.IsAvailable,
		dish.IsVegetarian,
		dish.IsFeatured,
		dish.PreparationTime,
	}

	return m.DB.QueryRow(query, args...).Scan(
		&dish.ID,
		&dish.CreatedAt,
		&dish.Version,
	)
}

func (m DishModel) Get(id int64) (*Dish, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, name, description, price,
		       category_id, image_url, is_available,
		       is_vegetarian, is_featured,
		       preparation_time, version
		FROM dishes
		WHERE id = $1`

	var dish Dish

	err := m.DB.QueryRow(query, id).Scan(
		&dish.ID,
		&dish.CreatedAt,
		&dish.Name,
		&dish.Description,
		&dish.Price,
		&dish.CategoryID,
		&dish.ImageURL,
		&dish.IsAvailable,
		&dish.IsVegetarian,
		&dish.IsFeatured,
		&dish.PreparationTime,
		&dish.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &dish, nil
}

func (m DishModel) Update(dish *Dish) error {
	query := `
		UPDATE dishes
		SET name = $1,
		    description = $2,
		    price = $3,
		    category_id = $4,
		    image_url = $5,
		    is_available = $6,
		    is_vegetarian = $7,
		    is_featured = $8,
		    preparation_time = $9,
		    version = version + 1
		WHERE id = $10
		AND version = $11
		RETURNING version`

	args := []interface{}{
		dish.Name,
		dish.Description,
		dish.Price,
		dish.CategoryID,
		dish.ImageURL,
		dish.IsAvailable,
		dish.IsVegetarian,
		dish.IsFeatured,
		dish.PreparationTime,
		dish.ID,
		dish.Version,
	}

	err := m.DB.QueryRow(query, args...).Scan(&dish.Version)

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

func (m DishModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM dishes
		WHERE id = $1`

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

func (m DishModel) ToggleAvailability(id int64, available bool) error {
	query := `
		UPDATE dishes
		SET is_available = $1,
		    version = version + 1
		WHERE id = $2`

	result, err := m.DB.Exec(query, available, id)
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

func (m DishModel) GetAll(categoryID int64, availableOnly bool) ([]*Dish, error) {
	query := `
		SELECT id, created_at, name, description,
		       price, category_id, image_url,
		       is_available, is_vegetarian,
		       is_featured, preparation_time,
		       version
		FROM dishes
		WHERE ($1 = 0 OR category_id = $1)
		  AND ($2 = false OR is_available = true)
		ORDER BY name ASC`

	rows, err := m.DB.Query(query, categoryID, availableOnly)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dishes := []*Dish{}

	for rows.Next() {
		var dish Dish

		err := rows.Scan(
			&dish.ID,
			&dish.CreatedAt,
			&dish.Name,
			&dish.Description,
			&dish.Price,
			&dish.CategoryID,
			&dish.ImageURL,
			&dish.IsAvailable,
			&dish.IsVegetarian,
			&dish.IsFeatured,
			&dish.PreparationTime,
			&dish.Version,
		)
		if err != nil {
			return nil, err
		}

		dishes = append(dishes, &dish)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return dishes, nil
}
