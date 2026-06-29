package data

import (
	"Dipu-36/restaurant/internals/validator"
	"database/sql"
	"errors"
	"time"
)

const (
	AnonymousUserID = 0
)

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Role      string    `json:"role"`
	Activated bool      `json:"activated"`
	Version   int32     `json:"version"`
	Password  password  `json:"-"`
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 100, "name", "must not be more than 100 bytes long")

	ValidateEmail(v, user.Email)

	v.Check(user.Phone != "", "phone", "must be provided")
	v.Check(len(user.Phone) >= 10, "phone", "must be a valid phone number")
	v.Check(len(user.Phone) <= 15, "phone", "must be a valid phone number")

	v.Check(user.Address != "", "address", "must be provided")
	v.Check(len(user.Address) <= 500, "address", "must not be more than 500 bytes long")

	v.Check(user.Role == "customer" || user.Role == "owner", "role", "must be customer or owner")

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash")
	}
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (name, email, password_hash, phone, address, role, activated)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id, created_at, version`

	args := []interface{}{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Phone,
		user.Address,
		user.Role,
		user.Activated,
	}

	return m.DB.QueryRow(query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Version,
	)
}

func (m UserModel) Get(id int64) (*User, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, name, email, password_hash,
		       phone, address, role, activated, version
		FROM users
		WHERE id = $1`

	var user User

	err := m.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Phone,
		&user.Address,
		&user.Role,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, created_at, name, email, password_hash,
		       phone, address, role, activated, version
		FROM users
		WHERE email = $1`

	var user User

	err := m.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Phone,
		&user.Address,
		&user.Role,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET name = $1,
		    email = $2,
		    password_hash = $3,
		    phone = $4,
		    address = $5,
		    role = $6,
		    activated = $7,
		    version = version + 1
		WHERE id = $8
		AND version = $9
		RETURNING version`

	args := []interface{}{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Phone,
		user.Address,
		user.Role,
		user.Activated,
		user.ID,
		user.Version,
	}
	err := m.DB.QueryRow(query, args...).Scan(&user.Version)

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

func (m UserModel) Delete(id int64) error {
	query := `
		DELETE FROM users
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
