package models

import "time"

type Order struct {
	ID                int64     `db:"id"`
	UserID            int64     `db:"user_id"`
	Status            string    `db:"status"`
	ShippingAddressID *int64    `db:"shipping_address_id"`
	TotalAmount       float64   `db:"total_amount"`
	Currency          string    `db:"currency"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

type OrderItem struct {
	ID        int64   `db:"id"`
	OrderID   int64   `db:"order_id"`
	ProductID int64   `db:"product_id"`
	Quantity  int     `db:"quantity"`
	Price     float64 `db:"price"`
}

type Payment struct {
	ID              int64     `db:"id"`
	OrderID         int64     `db:"order_id"`
	PaymentMethodID *int64    `db:"payment_method_id"`
	Status          string    `db:"status"`
	TransactionID   *string   `db:"transaction_id"`
	CreatedAt       time.Time `db:"created_at"`
}
