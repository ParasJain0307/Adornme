package migrations

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Migrator struct {
	Pool *pgxpool.Pool
}

func NewMigrator(pool *pgxpool.Pool) *Migrator {
	return &Migrator{Pool: pool}
}

// Dispatcher
func (m *Migrator) Migrate(ctx context.Context, dbName string) error {
	switch dbName {
	case "usersdb":
		return m.migrateUsers(ctx)
	case "productsdb":
		return m.migrateProducts(ctx)
	case "ordersdb":
		return m.migrateOrders(ctx)
	case "inventorydb":
		return m.migrateInventory(ctx)
	case "ecommerce":
		return m.migrateEcommerce(ctx)
	default:
		return fmt.Errorf("unknown dbName: %s", dbName)
	}
}

// ------------------ Users ------------------
func (m *Migrator) migrateUsers(ctx context.Context) error {
	_, err := m.Pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE,
		phone TEXT UNIQUE,
		password TEXT NOT NULL,
		email_verified BOOLEAN DEFAULT FALSE,
		phone_verified BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP,
		refresh_token TEXT,
		access_token TEXT
	);`)
	if err != nil {
		return err
	}
	if err := m.migratePasswordResets(ctx); err != nil {
		return err
	}

	return err
}

func (m *Migrator) migratePasswordResets(ctx context.Context) error {
	_, err := m.Pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS password_resets (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		token_hash TEXT NOT NULL,
		expiry TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_password_resets_token
	ON password_resets(token_hash);
	`)
	return err
}

// ------------------ Products ------------------
func (m *Migrator) migrateProducts(ctx context.Context) error {
	_, err := m.Pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		price NUMERIC(10,2) NOT NULL,
		inventory INT NOT NULL DEFAULT 0,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP
	);`)
	return err
}

// ------------------ Orders ------------------
func (m *Migrator) migrateOrders(ctx context.Context) error {
	_, err := m.Pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		total_amount NUMERIC(10,2) NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending',
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP
	);`)
	return err
}

// ------------------ Inventory ------------------
func (m *Migrator) migrateInventory(ctx context.Context) error {
	_, err := m.Pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS inventory (
		id SERIAL PRIMARY KEY,
		product_id INT NOT NULL,
		quantity INT NOT NULL DEFAULT 0,
		updated_at TIMESTAMP
	);`)
	return err
}

// ------------------ Ecommerce (extra tables) ------------------
func (m *Migrator) migrateEcommerce(ctx context.Context) error {
	_, err := m.Pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS cart_items (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		product_id INT NOT NULL,
		quantity INT NOT NULL DEFAULT 1,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`)
	return err
}
