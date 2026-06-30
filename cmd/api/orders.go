package main

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
