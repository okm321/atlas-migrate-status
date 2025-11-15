package db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Migration struct {
	Version       string    `json:"version"`
	Description   string    `json:"description"`
	ExecutedAt    time.Time `json:"executed_at"`
	ExecutionTime int64     `json:"execution_time"`
	Type          string    `json:"type"`
	Error         string    `json:"error,omitempty"`
}

func FetchMigrationHistory(dbURL string, tableName string) (migrations []Migration, err error) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer func() {
		if closeErr := conn.Close(ctx); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to close database connection: %w", closeErr))
		}
	}()

	query := fmt.Sprintf(`
		SELECT
			version,
			COALESCE(description, '') as description,
		  executed_at,
		  execution_time,
		  type,
		  COALESCE(error, '') as error
		FROM %s
		ORDER BY executed_at ASC
		`, tableName)

	rows, err := conn.Query(ctx, query)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil, fmt.Errorf("migration history table '%s' not found. Has 'atlas migrate apply' been run?", tableName)
		}
		return nil, fmt.Errorf("failed to query migration history: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m Migration
		err := rows.Scan(
			&m.Version,
			&m.Description,
			&m.ExecutedAt,
			&m.ExecutionTime,
			&m.Type,
			&m.Error,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		migrations = append(migrations, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return migrations, nil
}
