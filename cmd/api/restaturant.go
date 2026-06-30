package main

import (
	"Dipu-36/restaurant/internals/data"
	"Dipu-36/restaurant/internals/validator"
	"errors"
	"net/http"
)

type updateRestaurantInput struct {
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	Email          *string `json:"email"`
	Phone          *string `json:"phone"`
	StreetAddress  *string `json:"street_address"`
	OpeningTime    *string `json:"opening_time"`
	ClosingTime    *string `json:"closing_time"`
	DeliveryFee    *int64  `json:"delivery_fee"`
	DeliveryRadius *int32  `json:"delivery_radius"`
	IsOpen         *bool   `json:"is_open"`
}

func (app *application) getRestaurantHandler(w http.ResponseWriter, r *http.Request) {

	restaurant, err := app.models.Restaurants.Get()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(
		w,
		http.StatusOK,
		envelope{
			"restaurant": restaurant,
		},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateRestaurantHandler(w http.ResponseWriter, r *http.Request) {

	var input updateRestaurantInput

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	restaurant, err := app.models.Restaurants.Get()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if input.Name != nil {
		restaurant.Name = *input.Name
	}

	if input.Description != nil {
		restaurant.Description = *input.Description
	}

	if input.Email != nil {
		restaurant.Email = *input.Email
	}

	if input.Phone != nil {
		restaurant.Phone = *input.Phone
	}

	if input.StreetAddress != nil {
		restaurant.StreetAddress = *input.StreetAddress
	}

	if input.OpeningTime != nil {
		restaurant.OpeningTime = *input.OpeningTime
	}

	if input.ClosingTime != nil {
		restaurant.ClosingTime = *input.ClosingTime
	}

	if input.DeliveryFee != nil {
		restaurant.DeliveryFee = *input.DeliveryFee
	}

	if input.DeliveryRadius != nil {
		restaurant.DeliveryRadius = *input.DeliveryRadius
	}

	if input.IsOpen != nil {
		restaurant.IsOpen = *input.IsOpen
	}

	v := validator.New()

	if data.ValidateRestaurant(v, restaurant); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Restaurants.Update(restaurant)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(
		w,
		http.StatusOK,
		envelope{
			"restaurant": restaurant,
		},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
