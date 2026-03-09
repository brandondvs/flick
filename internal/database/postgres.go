package database

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/brandondvs/flick/internal/config"
)

type Postgres struct {
	conn *pgx.Conn
}

func (p *Postgres) Ping() error {
	return p.conn.Ping(context.Background())
}

func New() (*Postgres, error) {
	conn, err := pgx.Connect(context.Background(), config.DatabaseConnectionString())
	if err != nil {
		return nil, err
	}

	return &Postgres{
		conn,
	}, nil
}
