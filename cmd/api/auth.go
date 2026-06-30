package main

import (
	"errors"
	"net/http"

	"Dipu-36/restaurant/internals/data"
	"Dipu-36/restaurant/internals/validator"
)

type authenticationTokenInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {

	var input authenticationTokenInput

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	if !user.Activated {
		app.inactiveAccountResponse(w, r)
		return
	}

	authenticationToken, err := app.jwtManager.Generate(
		user.ID,
		user.Role,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(
		w,
		http.StatusCreated,
		envelope{
			"authentication_token": authenticationToken,
		},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
