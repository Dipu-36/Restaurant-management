package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	// -------------------------------------------------------------------------
	// Public Routes (no auth required)
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

	// Public menu browsing — customers and guests need these
	router.HandlerFunc(
		http.MethodGet,
		"/v1/categories",
		app.listCategoriesHandler,
	)

	router.HandlerFunc(
		http.MethodGet,
		"/v1/dishes",
		app.listDishesHandler,
	)

	router.HandlerFunc(
		http.MethodGet,
		"/v1/dishes/:id",
		app.getDishHandler,
	)

	router.HandlerFunc(
		http.MethodGet,
		"/v1/categories/:id",
		app.getCategoryHandler,
	)

	router.HandlerFunc(
		http.MethodGet,
		"/v1/restaurant",
		app.getRestaurantHandler,
	)

	// -------------------------------------------------------------------------
	// Owner Routes
	// -------------------------------------------------------------------------

	// Restaurant — only updates need owner protection
	router.Handler(
		http.MethodPatch,
		"/v1/restaurant",
		app.authenticate(
			app.requireOwner(
				http.HandlerFunc(app.updateRestaurantHandler),
			),
		),
	)

	// Categories — create, update, delete are owner-only
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
		http.MethodPatch,
		"/v1/categories/:id",
		app.authenticate(
			app.requireOwner(
				http.HandlerFunc(app.updateCategoryHandler),
			),
		),
	)

	router.Handler(
		http.MethodDelete,
		"/v1/categories/:id",
		app.authenticate(
			app.requireOwner(
				http.HandlerFunc(app.deleteCategoryHandler),
			),
		),
	)

	// Dishes — create, update, delete are owner-only
	router.Handler(
		http.MethodPost,
		"/v1/dishes",
		app.authenticate(
			app.requireOwner(
				http.HandlerFunc(app.createDish),
			),
		),
	)

	router.Handler(
		http.MethodPatch,
		"/v1/dishes/:id",
		app.authenticate(
			app.requireOwner(
				http.HandlerFunc(app.updateDishHandler),
			),
		),
	)

	router.Handler(
		http.MethodDelete,
		"/v1/dishes/:id",
		app.authenticate(
			app.requireOwner(
				http.HandlerFunc(app.deleteDishHandler),
			),
		),
	)

	// -------------------------------------------------------------------------
	// Customer Routes
	// -------------------------------------------------------------------------

	// Orders
	router.Handler(
		http.MethodPost,
		"/v1/orders",
		app.authenticate(
			app.requireCustomer(
				http.HandlerFunc(app.createOrderHandler),
			),
		),
	)

	router.Handler(
		http.MethodGet,
		"/v1/orders",
		app.authenticate(
			app.requireCustomer(
				http.HandlerFunc(app.listOrdersHandler),
			),
		),
	)

	router.Handler(
		http.MethodGet,
		"/v1/orders/:id",
		app.authenticate(
			app.requireCustomer(
				http.HandlerFunc(app.getOrderHandler),
			),
		),
	)

	router.Handler(
		http.MethodPatch,
		"/v1/orders/:id/cancel",
		app.authenticate(
			app.requireCustomer(
				http.HandlerFunc(app.cancelOrderHandler),
			),
		),
	)

	// Addresses
	router.Handler(
		http.MethodPost,
		"/v1/addresses",
		app.authenticate(
			app.requireCustomer(
				http.HandlerFunc(app.createAddressHandler),
			),
		),
	)

	router.Handler(
		http.MethodGet,
		"/v1/addresses",
		app.authenticate(
			app.requireCustomer(
				http.HandlerFunc(app.listAddressesHandler),
			),
		),
	)

	router.Handler(
		http.MethodGet,
		"/v1/addresses/:id",
		app.authenticate(
			app.requireCustomer(
				http.HandlerFunc(app.getAddressHandler),
			),
		),
	)

	router.Handler(
		http.MethodPatch,
		"/v1/addresses/:id",
		app.authenticate(
			app.requireCustomer(
				http.HandlerFunc(app.updateAddressHandler),
			),
		),
	)

	router.Handler(
		http.MethodDelete,
		"/v1/addresses/:id",
		app.authenticate(
			app.requireCustomer(
				http.HandlerFunc(app.deleteAddressHandler),
			),
		),
	)

	// -------------------------------------------------------------------------
	// Owner-Only Admin Routes
	// -------------------------------------------------------------------------

	router.Handler(
		http.MethodPatch,
		"/v1/orders/:id/status",
		app.authenticate(
			app.requireOwner(
				http.HandlerFunc(app.updateOrderStatusHandler),
			),
		),
	)

	router.Handler(
		http.MethodGet,
		"/v1/admin/orders",
		app.authenticate(
			app.requireOwner(
				http.HandlerFunc(app.listAllOrdersHandler),
			),
		),
	)

	router.Handler(
		http.MethodGet,
		"/v1/admin/dashboard",
		app.authenticate(
			app.requireOwner(
				http.HandlerFunc(app.getDashboardHandler),
			),
		),
	)

	return router
}