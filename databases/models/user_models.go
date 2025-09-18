package models

import "time"

type User struct {
	ID           int64     `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Phone        *string   `db:"phone"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Address struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Line1     string    `db:"line1"`
	Line2     *string   `db:"line2"`
	City      string    `db:"city"`
	State     string    `db:"state"`
	Country   string    `db:"country"`
	ZipCode   string    `db:"zipcode"`
	IsDefault bool      `db:"is_default"`
	CreatedAt time.Time `db:"created_at"`
}

type PaymentMethod struct {
	ID              int64     `db:"id"`
	UserID          int64     `db:"user_id"`
	CardLast4       string    `db:"card_last4"`
	CardType        *string   `db:"card_type"`
	ExpiryDate      *time.Time `db:"expiry_date"`
	BillingAddressID *int64   `db:"billing_address_id"`
	CreatedAt       time.Time `db:"created_at"`
}
