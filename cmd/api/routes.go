package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/createdish", app.createDish) // for the admin ONLY
	//router.HandlerFunc(http.MethodGet, "/", app.MenuHandler)
	//router.HandlerFunc(http.MethodGet, "/v1/menu", app.MenuHandler)

	return router
}
