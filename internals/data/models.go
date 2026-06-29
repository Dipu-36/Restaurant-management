package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Dishes     DishModel
	Users      UserModel
	Categories CategoryModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Dishes:     DishModel{DB: db},
		Users:      UserModel{DB: db},
		Categories: CategoryModel{DB: db},
	}
}
