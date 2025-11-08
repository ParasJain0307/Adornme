package database

import (
	migrations "Adornme/databases/migration"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ---------------- Postgres ----------------
type PostgresProvider struct {
	Pool *pgxpool.Pool
}

// Postgres component clients
type PostgresClients struct {
	UsersDB     *PostgresProvider
	ProductsDB  *PostgresProvider
	OrdersDB    *PostgresProvider
	InventoryDB *PostgresProvider
}

// Implement DatabaseProvider for PostgresClients
func (c *PostgresClients) HealthCheck(ctx context.Context) error {
	for _, db := range []*PostgresProvider{c.UsersDB, c.ProductsDB, c.OrdersDB, c.InventoryDB} {
		if db != nil {
			if err := db.HealthCheck(ctx); err != nil {
				return fmt.Errorf("postgres client unhealthy: %w", err)
			}
		}
	}
	return nil
}

func (c *PostgresClients) Close() error {
	for _, db := range []*PostgresProvider{c.UsersDB, c.ProductsDB, c.OrdersDB, c.InventoryDB} {
		if db != nil {
			db.Close()
		}
	}
	return nil
}

// Individual PostgresProvider methods
func (p *PostgresProvider) HealthCheck(ctx context.Context) error {
	return p.Pool.Ping(ctx)
}

func (p *PostgresProvider) Close() error {
	if p.Pool != nil {
		p.Pool.Close()
	}
	return nil
}

// Postgres config
type postgresConfig struct {
	Enabled         bool   `json:"enabled"`
	DSN             string `json:"dsn"`
	MaxConns        int32  `json:"maxConns"`
	ConnMaxLifetime string `json:"connMaxLifetime"`
}

func connectPostgres(raw json.RawMessage) (*PostgresProvider, error) {
	var cfg postgresConfig

	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return nil, nil
	}

	// Parse DSN
	u, err := pgx.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, err
	}
	dbName := u.Database

	// Connect to the default "postgres" database to create target DB if it doesn't exist
	u.Database = "postgres" // default DB exists
	defaultDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=disable",
		u.User, u.Password, u.Host, u.Port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, defaultDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to default DB: %w", err)
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, dbName))
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return nil, fmt.Errorf("failed to create database %s: %w", dbName, err)
	}

	// Now connect to the actual target database
	pool, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		return nil, err
	}

	logs.Infof(Ctx, "Postgres connected ✅ for %v", dbName)
	// run migrations
	m := migrations.NewMigrator(pool)
	if err := m.Migrate(ctx, dbName); err != nil {
		logs.Errorf(ctx, "migration failed:%v", err)
	}
	logs.Infof(Ctx, "migration successful for db:%v", dbName)
	return &PostgresProvider{Pool: pool}, nil
}

// Connect all Postgres components
func ConnectAllPostgres(configs map[string]json.RawMessage) (*PostgresClients, error) {
	clients := &PostgresClients{}
	var err error

	if r, ok := configs["users"]; ok {
		clients.UsersDB, err = connectPostgres(r)
		if err != nil {
			return nil, err
		}
	}
	if r, ok := configs["products"]; ok {
		clients.ProductsDB, err = connectPostgres(r)
		if err != nil {
			return nil, err
		}
	}
	if r, ok := configs["orders"]; ok {
		clients.OrdersDB, err = connectPostgres(r)
		if err != nil {
			return nil, err
		}
	}
	if r, ok := configs["inventory"]; ok {
		clients.InventoryDB, err = connectPostgres(r)
		if err != nil {
			return nil, err
		}
	}
	logs.Info(Ctx, "All Postgres Connection Made Successfully")
	return clients, nil
}
