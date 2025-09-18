package database

import (
	"context"
	"time"
)

// ----------------- User Model -----------------
type User struct {
	ID        int       `db:"id"`         // Primary Key
	Name      string    `db:"name"`       // User full name
	Email     string    `db:"email"`      // Unique email
	Password  string    `db:"password"`   // Hashed password
	CreatedAt time.Time `db:"created_at"` // Creation timestamp
	UpdatedAt time.Time `db:"updated_at"` // Optional update timestamp
}

// ----------------- Product Model -----------------
type Product struct {
	ID          int       `db:"id"`          // Primary Key
	Name        string    `db:"name"`        // Product name
	Description string    `db:"description"` // Product description
	Price       float64   `db:"price"`       // Product price
	Inventory   int       `db:"inventory"`   // Stock quantity
	CreatedAt   time.Time `db:"created_at"`  // Creation timestamp
	UpdatedAt   time.Time `db:"updated_at"`  // Optional update timestamp
}

// ----------------- Order Model -----------------
type Order struct {
	ID        int       `db:"id"`         // Primary Key
	UserID    int       `db:"user_id"`    // Foreign key to users
	Total     float64   `db:"total"`      // Total order amount
	Status    string    `db:"status"`     // Pending, Completed, Cancelled
	CreatedAt time.Time `db:"created_at"` // Creation timestamp
	UpdatedAt time.Time `db:"updated_at"` // Optional update timestamp
}

// ----------------- OrderItem Model -----------------
type OrderItem struct {
	ID        int     `db:"id"`         // Primary Key
	OrderID   int     `db:"order_id"`   // Foreign key to orders
	ProductID int     `db:"product_id"` // Foreign key to products
	Quantity  int     `db:"quantity"`   // Quantity of this product
	Price     float64 `db:"price"`      // Price per unit
}

// ----------------- Inventory Model -----------------
type Inventory struct {
	ID        int       `db:"id"`         // Primary Key
	ProductID int       `db:"product_id"` // Foreign key to products
	Quantity  int       `db:"quantity"`   // Current stock
	UpdatedAt time.Time `db:"updated_at"` // Last stock update
}

// ----------------- User CRUD -----------------
func (p *PostgresProvider) CreateUser(ctx context.Context, u User) (int, error) {
	logs.Info(ctx, "Created User")
	var id int
	err := p.Pool.QueryRow(ctx,
		`INSERT INTO users (name,email,password,created_at) VALUES ($1,$2,$3,$4) RETURNING id`,
		u.Name, u.Email, u.Password, u.CreatedAt).Scan(&id)
	return id, err
}

func (p *PostgresProvider) GetUser(ctx context.Context, id int) (*User, error) {
	u := &User{}
	err := p.Pool.QueryRow(ctx,
		`SELECT id,name,email,password,created_at FROM users WHERE id=$1`, id).
		Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (p *PostgresProvider) UpdateUser(ctx context.Context, u User) error {
	_, err := p.Pool.Exec(ctx,
		`UPDATE users SET name=$1,email=$2,password=$3 WHERE id=$4`,
		u.Name, u.Email, u.Password, u.ID)
	return err
}

func (p *PostgresProvider) DeleteUser(ctx context.Context, id int) error {
	_, err := p.Pool.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	return err
}

// ----------------- Product CRUD -----------------
func (p *PostgresProvider) CreateProduct(ctx context.Context, prod Product) (int, error) {
	var id int
	err := p.Pool.QueryRow(ctx,
		`INSERT INTO products (name,description,price,inventory,created_at) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		prod.Name, prod.Description, prod.Price, prod.Inventory, prod.CreatedAt).Scan(&id)
	return id, err
}

func (p *PostgresProvider) GetProduct(ctx context.Context, id int) (*Product, error) {
	prod := &Product{}
	err := p.Pool.QueryRow(ctx,
		`SELECT id,name,description,price,inventory,created_at FROM products WHERE id=$1`, id).
		Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.Inventory, &prod.CreatedAt)
	return prod, err
}

func (p *PostgresProvider) UpdateProduct(ctx context.Context, prod Product) error {
	_, err := p.Pool.Exec(ctx,
		`UPDATE products SET name=$1,description=$2,price=$3,inventory=$4 WHERE id=$5`,
		prod.Name, prod.Description, prod.Price, prod.Inventory, prod.ID)
	return err
}

func (p *PostgresProvider) DeleteProduct(ctx context.Context, id int) error {
	_, err := p.Pool.Exec(ctx, `DELETE FROM products WHERE id=$1`, id)
	return err
}

// ----------------- Order CRUD -----------------
func (p *PostgresProvider) CreateOrder(ctx context.Context, order Order, items []OrderItem) (int, error) {
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var orderID int
	err = tx.QueryRow(ctx,
		`INSERT INTO orders (user_id,total,status,created_at) VALUES ($1,$2,$3,$4) RETURNING id`,
		order.UserID, order.Total, order.Status, order.CreatedAt).Scan(&orderID)
	if err != nil {
		return 0, err
	}

	for _, item := range items {
		_, err := tx.Exec(ctx,
			`INSERT INTO order_items (order_id,product_id,quantity,price) VALUES ($1,$2,$3,$4)`,
			orderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return orderID, nil
}

func (p *PostgresProvider) GetOrder(ctx context.Context, orderID int) (*Order, []OrderItem, error) {
	order := &Order{}
	err := p.Pool.QueryRow(ctx,
		`SELECT id,user_id,total,status,created_at FROM orders WHERE id=$1`, orderID).
		Scan(&order.ID, &order.UserID, &order.Total, &order.Status, &order.CreatedAt)
	if err != nil {
		return nil, nil, err
	}

	rows, err := p.Pool.Query(ctx,
		`SELECT id,order_id,product_id,quantity,price FROM order_items WHERE order_id=$1`, orderID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	items := []OrderItem{}
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
			return nil, nil, err
		}
		items = append(items, item)
	}

	return order, items, nil
}
