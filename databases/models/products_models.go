package models

import "time"

type Category struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	ParentID *int64 `db:"parent_id"`
}

type Product struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`
	CategoryID  *int64    `db:"category_id"`
	Brand       *string   `db:"brand"`
	Price       float64   `db:"price"`
	Currency    string    `db:"currency"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type ProductImage struct {
	ID        int64  `db:"id"`
	ProductID int64  `db:"product_id"`
	ImageURL  string `db:"image_url"`
}

type ProductAttribute struct {
	ID             int64  `db:"id"`
	ProductID      int64  `db:"product_id"`
	AttributeName  string `db:"attribute_name"`
	AttributeValue string `db:"attribute_value"`
}
