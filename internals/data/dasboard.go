package data

type DashboardStats struct {
	RestaurantName string `json:"restaurant_name"`

	CategoryCount int `json:"category_count"`
	DishCount     int `json:"dish_count"`

	PendingOrders        int `json:"pending_orders"`
	AcceptedOrders       int `json:"accepted_orders"`
	PreparingOrders      int `json:"preparing_orders"`
	ReadyOrders          int `json:"ready_orders"`
	OutForDeliveryOrders int `json:"out_for_delivery_orders"`
	DeliveredOrders      int `json:"delivered_orders"`
	CancelledOrders      int `json:"cancelled_orders"`
}

type DashboardModel struct {
	DB DBTX
}

func (m DashboardModel) Get() (*DashboardStats, error) {

	stats := &DashboardStats{}

	err := m.DB.QueryRow(`
		SELECT name
		FROM restaurant
		LIMIT 1
	`).Scan(&stats.RestaurantName)

	if err != nil {
		return nil, err
	}

	err = m.DB.QueryRow(`
		SELECT COUNT(*)
		FROM categories
	`).Scan(&stats.CategoryCount)

	if err != nil {
		return nil, err
	}

	err = m.DB.QueryRow(`
		SELECT COUNT(*)
		FROM dishes
	`).Scan(&stats.DishCount)

	if err != nil {
		return nil, err
	}

	rows, err := m.DB.Query(`
		SELECT status, COUNT(*)
		FROM orders
		GROUP BY status
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var status string
		var count int

		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}

		switch status {

		case OrderStatusPending:
			stats.PendingOrders = count

		case OrderStatusAccepted:
			stats.AcceptedOrders = count

		case OrderStatusPreparing:
			stats.PreparingOrders = count

		case OrderStatusReady:
			stats.ReadyOrders = count

		case OrderStatusOutForDelivery:
			stats.OutForDeliveryOrders = count

		case OrderStatusDelivered:
			stats.DeliveredOrders = count

		case OrderStatusCancelled:
			stats.CancelledOrders = count
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}
