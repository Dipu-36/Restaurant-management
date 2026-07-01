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

	// Restaurant
	router.Handler(
		http.MethodGet,
		"/v1/restaurant",
		app.authenticate(
			app.requireOwner(
				http.HandlerFunc(app.getRestaurantHandler),
			),
		),
	)

	router.Handler(
		http.MethodPatch,
		"/v1/restaurant",
		app.authenticate(
			app.requireOwner(
				http.HandlerFunc(app.updateRestaurantHandler),
			),
		),
	)

	// Categories
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
		http.MethodGet,
		"/v1/categories",
		app.authenticate(
			app.requireOwner(
				http.HandlerFunc(app.listCategoriesHandler),
			),
		),
	)

	router.Handler(
		http.MethodGet,
		"/v1/categories/:id",
		app.authenticate(
			app.requireOwner(
				http.HandlerFunc(app.getCategoryHandler),
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

	// Dishes
	// -------------------------------------------------------------------------
	// Customer Routes
	// -------------------------------------------------------------------------

	router.Handler(
		http.MethodPost,
		"/v1/orders",
		app.authenticate(
			app.requireCustomer(
				http.HandlerFunc(app.createOrderHandler),
			),
		),
	)
	// -------------------------------------------------------------------------
	// Dishes
	// -------------------------------------------------------------------------

	router.Handler(
		http.MethodGet,
		"/v1/dishes",
		app.authenticate(
			app.requireOwner(
				http.HandlerFunc(app.listDishesHandler),
			),
		),
	)

	router.Handler(
		http.MethodGet,
		"/v1/dishes/:id",
		app.authenticate(
			app.requireOwner(
				http.HandlerFunc(app.getDishHandler),
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
	// Addresses
	// -------------------------------------------------------------------------

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
		http.MethodGet,
		"/v1/orders",
		app.authenticate(
			app.requireCustomer(
				http.HandlerFunc(app.listOrdersHandler),
			),
		),
	)
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
		http.MethodPatch,
		"/v1/orders/:id/cancel",
		app.authenticate(
			app.requireCustomer(
				http.HandlerFunc(app.cancelOrderHandler),
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
