package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	// Health Check
	router.HandlerFunc(
		http.MethodGet,
		"/v1/healthcheck",
		app.healthCheckHandler,
	)

	// Authentication
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

	// Restaurant
	router.Handler(
		http.MethodPost,
		"/v1/dishes",
		app.authenticate(
			http.HandlerFunc(app.createDish),
		),
	)

	// Future Routes
	// router.Handler(http.MethodGet, "/v1/restaurants", app.authenticate(http.HandlerFunc(app.listRestaurants)))
	// router.Handler(http.MethodPost, "/v1/orders", app.authenticate(http.HandlerFunc(app.createOrder)))
	// router.Handler(http.MethodGet, "/v1/orders/:id", app.authenticate(http.HandlerFunc(app.getOrder)))
	// router.Handler(http.MethodPost, "/v1/payments", app.authenticate(http.HandlerFunc(app.createPayment)))

	return router
}
