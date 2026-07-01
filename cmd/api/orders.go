package main

import (
	"Dipu-36/restaurant/internals/data"
	"Dipu-36/restaurant/internals/validator"
	"errors"
	"net/http"
	"time"
)

type createOrderItemInput struct {
	DishID   int64 `json:"dish_id"`
	Quantity int32 `json:"quantity"`
}

type createOrderInput struct {
	OrderType         string                 `json:"order_type"`
	DeliveryAddressID *int64                 `json:"delivery_address_id,omitempty"`
	PickupTime        *string                `json:"pickup_time,omitempty"`
	Items             []createOrderItemInput `json:"items"`
}

func (app *application) createOrderHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	var input createOrderInput

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if len(input.Items) == 0 {
		app.failedValidationResponse(
			w,
			r,
			map[string]string{
				"items": "must contain at least one item",
			},
		)
		return
	}

	order := &data.Order{
		CustomerID:  user.ID,
		OrderStatus: data.OrderStatusPending,
		OrderType:   input.OrderType,
	}

	if input.OrderType == data.OrderTypeDelivery {

		if input.DeliveryAddressID == nil {
			app.failedValidationResponse(
				w,
				r,
				map[string]string{
					"delivery_address_id": "must be provided",
				},
			)
			return
		}

		order.DeliveryAddressID = *input.DeliveryAddressID

	} else if input.OrderType == data.OrderTypePickup {

		if input.PickupTime == nil {
			app.failedValidationResponse(
				w,
				r,
				map[string]string{
					"pickup_time": "must be provided",
				},
			)
			return
		}

		pickupTime, err := time.Parse(time.RFC3339, *input.PickupTime)
		if err != nil {
			app.failedValidationResponse(
				w,
				r,
				map[string]string{
					"pickup_time": "must be a valid RFC3339 timestamp",
				},
			)
			return
		}

		order.PickupTime = pickupTime
	}

	if order.OrderType == data.OrderTypeDelivery {

		address, err := app.models.Addresses.Get(order.DeliveryAddressID)
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
	}

	var orderItems []*data.OrderItem

	for _, inputItem := range input.Items {

		dish, err := app.models.Dishes.Get(inputItem.DishID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		if !dish.IsAvailable {
			app.failedValidationResponse(
				w,
				r,
				map[string]string{
					"dish": dish.Name + " is currently unavailable",
				},
			)
			return
		}

		item := &data.OrderItem{
			DishID:          dish.ID,
			Quantity:        inputItem.Quantity,
			PriceAtPurchase: dish.Price,
		}

		item.CalculateSubtotal()

		v := validator.New()

		data.ValidateOrderItem(v, item)

		if !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		orderItems = append(orderItems, item)

		order.Subtotal += item.Subtotal
	}

	order.Tax = order.Subtotal * data.DefaultTaxPercentage / 100

	if order.OrderType == data.OrderTypeDelivery {
		order.DeliveryFee = data.DefaultDeliveryFee
	} else {
		order.DeliveryFee = 0
	}

	order.CalculateTotal()

	v := validator.New()

	data.ValidateOrder(v, order)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	tx, err := app.db.Begin()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	defer tx.Rollback()

	models := data.NewModels(tx)

	err = models.Orders.Insert(order)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	for _, item := range orderItems {

		item.OrderID = order.ID

		err = models.OrderItems.Insert(item)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(
		w,
		http.StatusCreated,
		envelope{
			"order": order,
			"items": orderItems,
		},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getOrderHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)

	order, err := app.models.Orders.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Customers can only view their own orders.
	if order.CustomerID != user.ID {
		app.notPermittedResponse(w, r)
		return
	}

	items, err := app.models.OrderItems.GetByOrderID(order.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(
		w,
		http.StatusOK,
		envelope{
			"order": order,
			"items": items,
		},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listOrdersHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	orders, err := app.models.Orders.GetCustomerOrders(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(
		w,
		http.StatusOK,
		envelope{
			"orders": orders,
		},
		nil,
	)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

type updateOrderStatusInput struct {
	Status string `json:"status"`
}

func (app *application) updateOrderStatusHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	order, err := app.models.Orders.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input updateOrderStatusInput

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if !order.CanTransitionTo(input.Status) {
		app.failedValidationResponse(
			w,
			r,
			map[string]string{
				"status": "invalid status transition",
			},
		)
		return
	}

	order.OrderStatus = input.Status

	err = app.models.Orders.UpdateStatus(order)
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
			"order": order,
		},
		nil,
	)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) cancelOrderHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)

	order, err := app.models.Orders.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Customers can only cancel their own orders.
	if order.CustomerID != user.ID {
		app.notPermittedResponse(w, r)
		return
	}

	// Only pending or accepted orders can be cancelled.
	if order.OrderStatus != data.OrderStatusPending &&
		order.OrderStatus != data.OrderStatusAccepted {

		app.failedValidationResponse(
			w,
			r,
			map[string]string{
				"status": "this order can no longer be cancelled",
			},
		)
		return
	}

	order.OrderStatus = data.OrderStatusCancelled

	err = app.models.Orders.UpdateStatus(order)
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
			"order": order,
		},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listAllOrdersHandler(w http.ResponseWriter, r *http.Request) {

	status := r.URL.Query().Get("status")

	orders, err := app.models.Orders.GetAll(status)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(
		w,
		http.StatusOK,
		envelope{
			"orders": orders,
		},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
