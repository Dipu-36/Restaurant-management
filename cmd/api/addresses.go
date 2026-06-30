package main

import (
	"Dipu-36/restaurant/internals/data"
	"Dipu-36/restaurant/internals/validator"
	"errors"
	"net/http"
)

type createAddressInput struct {
	StreetLine1 string `json:"street_line_1"`
	StreetLine2 string `json:"street_line_2"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postal_code"`
	Country     string `json:"country"`
	IsDefault   bool   `json:"is_default"`
}

type updateAddressInput struct {
	StreetLine1 *string `json:"street_line_1"`
	StreetLine2 *string `json:"street_line_2"`
	City        *string `json:"city"`
	State       *string `json:"state"`
	PostalCode  *string `json:"postal_code"`
	Country     *string `json:"country"`
	IsDefault   *bool   `json:"is_default"`
}

func (app *application) createAddressHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	var input createAddressInput

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	address := &data.Address{
		CustomerID:  user.ID,
		StreetLine1: input.StreetLine1,
		StreetLine2: input.StreetLine2,
		City:        input.City,
		State:       input.State,
		PostalCode:  input.PostalCode,
		Country:     input.Country,
		IsDefault:   input.IsDefault,
	}

	v := validator.New()

	if data.ValidateAddress(v, address); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Addresses.Insert(address)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(
		w,
		http.StatusCreated,
		envelope{
			"address": address,
		},
		nil,
	)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAddressHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	address, err := app.models.Addresses.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if address.CustomerID != user.ID {
		app.notPermittedResponse(w, r)
		return
	}

	err = app.writeJSON(
		w,
		http.StatusOK,
		envelope{
			"address": address,
		},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listAddressesHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	addresses, err := app.models.Addresses.GetAll(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(
		w,
		http.StatusOK,
		envelope{
			"addresses": addresses,
		},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateAddressHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	address, err := app.models.Addresses.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if address.CustomerID != user.ID {
		app.notPermittedResponse(w, r)
		return
	}

	var input updateAddressInput

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.StreetLine1 != nil {
		address.StreetLine1 = *input.StreetLine1
	}

	if input.StreetLine2 != nil {
		address.StreetLine2 = *input.StreetLine2
	}

	if input.City != nil {
		address.City = *input.City
	}

	if input.State != nil {
		address.State = *input.State
	}

	if input.PostalCode != nil {
		address.PostalCode = *input.PostalCode
	}

	if input.Country != nil {
		address.Country = *input.Country
	}

	if input.IsDefault != nil {
		address.IsDefault = *input.IsDefault
	}

	v := validator.New()

	if data.ValidateAddress(v, address); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Addresses.Update(address)
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
			"address": address,
		},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteAddressHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	address, err := app.models.Addresses.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if address.CustomerID != user.ID {
		app.notPermittedResponse(w, r)
		return
	}

	err = app.models.Addresses.Delete(id)
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
			"message": "address deleted successfully",
		},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
