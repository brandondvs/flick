package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/brandondvs/flick/internal/config"
)

type Postgres struct {
	conn *pgx.Conn
}

func (p *Postgres) Ping() error {
	return p.conn.Ping(context.Background())
}

func (p *Postgres) Query(sql string, args ...any) ([]any, error) {
	rows, err := p.conn.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, errors.New("rows next is false after running query")
	}

	vals, err := rows.Values()
	if err != nil {
		return nil, err
	}

	return vals, nil
}

func ScanValue[T any](rows []any, index int) (T, error) {
	var zero T
	if index >= len(rows) {
		return zero, fmt.Errorf("index %d is out of range (len %d)", len(rows))
	}

	val, ok := rows[index].(T)
	if !ok {
		return zero, fmt.Errorf("expect %T at index %d, got %T", zero, index, rows[index])
	}
	return val, nil
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
