package models

import "time"

type Warehouse struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Location string `db:"location"`
}

type Inventory struct {
	ID          int64 `db:"id"`
	ProductID   int64 `db:"product_id"`
	WarehouseID int64 `db:"warehouse_id"`
	Quantity    int   `db:"quantity"`
}

type Supplier struct {
	ID          int64  `db:"id"`
	Name        string `db:"name"`
	ContactInfo string `db:"contact_info"`
}

type PurchaseOrder struct {
	ID         int64     `db:"id"`
	SupplierID int64     `db:"supplier_id"`
	ProductID  int64     `db:"product_id"`
	Quantity   int       `db:"quantity"`
	Status     string    `db:"status"`
	OrderedAt  time.Time `db:"ordered_at"`
}
