package data

import (
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrDuplicateEmail = errors.New("duplicate email")
)

type Models struct {
	Users       UserModel
	Restaurants RestaurantModel
	Categories  CategoryModel
	Dishes      DishModel
	Addresses   AddressModel
	Orders      OrderModel
	OrderItems  OrderItemModel
	Payments    PaymentModel
	Tokens      TokenModel
}

func NewModels(db DBTX) Models {
	return Models{
		Users:       UserModel{DB: db},
		Restaurants: RestaurantModel{DB: db},
		Categories:  CategoryModel{DB: db},
		Dishes:      DishModel{DB: db},
		Addresses:   AddressModel{DB: db},
		Orders:      OrderModel{DB: db},
		OrderItems:  OrderItemModel{DB: db},
		Payments:    PaymentModel{DB: db},
		Tokens:      TokenModel{DB: db},
	}
}
