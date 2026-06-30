package main

import (
	"errors"
	"net/http"
	"strconv"

	"Dipu-36/restaurant/internals/data"
	"Dipu-36/restaurant/internals/validator"
)

type createDishInput struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	Price           int64  `json:"price"`
	CategoryID      int64  `json:"category_id"`
	ImageURL        string `json:"image_url"`
	IsAvailable     bool   `json:"is_available"`
	IsVegetarian    bool   `json:"is_vegetarian"`
	IsFeatured      bool   `json:"is_featured"`
	PreparationTime int32  `json:"preparation_time"`
}

type updateDishInput struct {
	Name            *string `json:"name"`
	Description     *string `json:"description"`
	Price           *int64  `json:"price"`
	CategoryID      *int64  `json:"category_id"`
	ImageURL        *string `json:"image_url"`
	IsAvailable     *bool   `json:"is_available"`
	IsVegetarian    *bool   `json:"is_vegetarian"`
	IsFeatured      *bool   `json:"is_featured"`
	PreparationTime *int32  `json:"preparation_time"`
}

func (app *application) createDish(w http.ResponseWriter, r *http.Request) {

	var input createDishInput

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	dish := &data.Dish{
		Name:            input.Name,
		Description:     input.Description,
		Price:           input.Price,
		CategoryID:      input.CategoryID,
		ImageURL:        input.ImageURL,
		IsAvailable:     input.IsAvailable,
		IsVegetarian:    input.IsVegetarian,
		IsFeatured:      input.IsFeatured,
		PreparationTime: input.PreparationTime,
	}

	v := validator.New()

	if data.ValidateDish(v, dish); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Dishes.Insert(dish)
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
		http.StatusCreated,
		envelope{
			"dish": dish,
		},
		nil,
	)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateDishHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	dish, err := app.models.Dishes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input updateDishInput

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		dish.Name = *input.Name
	}

	if input.Description != nil {
		dish.Description = *input.Description
	}

	if input.Price != nil {
		dish.Price = *input.Price
	}

	if input.CategoryID != nil {
		dish.CategoryID = *input.CategoryID
	}

	if input.ImageURL != nil {
		dish.ImageURL = *input.ImageURL
	}

	if input.IsAvailable != nil {
		dish.IsAvailable = *input.IsAvailable
	}

	if input.IsVegetarian != nil {
		dish.IsVegetarian = *input.IsVegetarian
	}

	if input.IsFeatured != nil {
		dish.IsFeatured = *input.IsFeatured
	}

	if input.PreparationTime != nil {
		dish.PreparationTime = *input.PreparationTime
	}

	v := validator.New()

	if data.ValidateDish(v, dish); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Dishes.Update(dish)
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
			"dish": dish,
		},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteDishHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Dishes.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(
		w,
		http.StatusOK,
		envelope{
			"message": "dish deleted successfully",
		},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) getDishHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	dish, err := app.models.Dishes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(
		w,
		http.StatusOK,
		envelope{
			"dish": dish,
		},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listDishesHandler(w http.ResponseWriter, r *http.Request) {

	var (
		categoryID    int64
		availableOnly bool
	)

	query := r.URL.Query()

	if category := query.Get("category_id"); category != "" {
		id, err := strconv.ParseInt(category, 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, errors.New("invalid category_id"))
			return
		}
		categoryID = id
	}

	if available := query.Get("available_only"); available != "" {
		value, err := strconv.ParseBool(available)
		if err != nil {
			app.badRequestResponse(w, r, errors.New("invalid available_only"))
			return
		}
		availableOnly = value
	}

	dishes, err := app.models.Dishes.GetAll(categoryID, availableOnly)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(
		w,
		http.StatusOK,
		envelope{
			"dishes": dishes,
		},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
