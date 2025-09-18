package models

import "time"

type EventLog struct {
	ID          int64     `db:"id"`
	ServiceName string    `db:"service_name"`
	EntityID    int64     `db:"entity_id"`
	Action      string    `db:"action"`
	Payload     []byte    `db:"payload"` // JSONB
	CreatedAt   time.Time `db:"created_at"`
}

type UserActivity struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Action    string    `db:"action"`
	Metadata  []byte    `db:"metadata"` // JSONB
	CreatedAt time.Time `db:"created_at"`
}

type SalesAnalytics struct {
	ID           int64     `db:"id"`
	Date         time.Time `db:"date"`
	TotalOrders  int       `db:"total_orders"`
	TotalRevenue float64   `db:"total_revenue"`
	TopProductID *int64    `db:"top_product_id"`
	CreatedAt    time.Time `db:"created_at"`
}
