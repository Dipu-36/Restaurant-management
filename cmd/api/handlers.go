package main

import (
	"Dipu-36/restaurant/internals/validator"
	"net/http"
)

func (app *application) createDish(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string
		Type string
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	v.Check(input.Name != "", "name", "must be provided")
	v.Check(len(input.Name) <= 100, "name", "must not be more than 100 bytes long")
}
