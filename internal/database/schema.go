package database

import (
	"fmt"

	"github.com/brandondvs/flick/internal/config"
)

type SchemaBuilder struct {
	db *Postgres
}

func (s *SchemaBuilder) IsValid() error {
	rows, err := s.db.Query("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", config.DatabaseName())
	if err != nil {
		return err
	}

	dbExists, err := ScanValue[bool](rows, 0)
	if err != nil {
		return err
	}
	if !dbExists {
		return fmt.Errorf("database %s does not exist", config.DatabaseName())
	}
	return nil
}

func Schema(db *Postgres) *SchemaBuilder {
	return &SchemaBuilder{
		db,
	}
}
