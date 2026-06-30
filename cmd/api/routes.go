package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	// -------------------------------------------------------------------------
	// Public Routes
	// -------------------------------------------------------------------------

	router.HandlerFunc(
		http.MethodGet,
		"/v1/healthcheck",
		app.healthCheckHandler,
	)

	router.HandlerFunc(
		http.MethodPost,
		"/v1/users",
		app.registerUserHandler,
	)

	router.HandlerFunc(
		http.MethodPut,
		"/v1/users/activated",
		app.activateUserHandler,
	)

	router.HandlerFunc(
		http.MethodPost,
		"/v1/tokens/authentication",
		app.createAuthenticationTokenHandler,
	)

	// -------------------------------------------------------------------------
	// Owner Routes
	// -------------------------------------------------------------------------

	router.Handler(
		http.MethodPost,
		"/v1/categories",
		app.authenticate(
			app.requireOwner(
				http.HandlerFunc(app.createCategory),
			),
		),
	)

	router.Handler(
		http.MethodPost,
		"/v1/dishes",
		app.authenticate(
			app.requireOwner(
				http.HandlerFunc(app.createDish),
			),
		),
	)

	// -------------------------------------------------------------------------
	// Customer Routes
	// -------------------------------------------------------------------------

	router.Handler(
		http.MethodPost,
		"/v1/orders",
		app.authenticate(
			app.requireCustomer(
				http.HandlerFunc(app.createOrder),
			),
		),
	)

	return router
}
