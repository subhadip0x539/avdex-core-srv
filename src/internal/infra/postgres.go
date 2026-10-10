package infra

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool     *pgxpool.Pool
	host     string
	port     int32
	database string
	username string
	password string
}

func (d *Postgres) Host(host string) *Postgres {
	d.host = host
	return d
}

func (d *Postgres) Port(port int32) *Postgres {
	d.port = port
	return d
}

func (d *Postgres) Database(database string) *Postgres {
	d.database = database
	return d
}

func (d *Postgres) Username(username string) *Postgres {
	d.username = username
	return d
}

func (d *Postgres) Password(password string) *Postgres {
	d.password = password
	return d
}

func (d *Postgres) Build(ctx context.Context) (*Postgres, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", d.username, d.password, d.host, d.port, d.database)

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse db config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	d.pool = pool
	return d, nil
}

func (d *Postgres) Pool() *pgxpool.Pool {
	return d.pool
}

func (d *Postgres) Close() {
	if d.pool != nil {
		d.pool.Close()
	}
}

func NewPostgres() *Postgres {
	return &Postgres{}
}
