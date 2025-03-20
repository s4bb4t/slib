package migrator

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
)

type Connector interface {
	DB() *sql.DB
}

// Migrate performs database migrations using the specified connection, table name, and SQL dialect.
// Migrate use dir named `migrations` in root of project
func Migrate(conn Connector, tableName, dialect string) error {
	goose.SetTableName(fmt.Sprintf("public.%s", tableName))

	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("dialect error: %w", err)
	}

	if err := goose.Up(conn.DB(), "migrations"); err != nil {
		return fmt.Errorf("migrate error: %w", err)
	}

	return nil
}

// Downgrade rolls back the latest database migration for the specified table using the provided connector and dialect.
// Downgrade use dir named `migrations` in root of project
func Downgrade(conn Connector, tableName, dialect string) error {
	goose.SetTableName(fmt.Sprintf("public.%s", tableName))

	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("dialect error: %w", err)
	}

	if err := goose.Down(conn.DB(), "migrations/"); err != nil {
		return fmt.Errorf("migrate error: %w", err)
	}

	return nil
}
