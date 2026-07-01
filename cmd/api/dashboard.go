package main

import (
	"net/http"
)

func (app *application) getDashboardHandler(w http.ResponseWriter, r *http.Request) {

	stats, err := app.models.Dashboard.Get()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(
		w,
		http.StatusOK,
		envelope{
			"dashboard": stats,
		},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
