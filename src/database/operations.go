package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// CreateSchemaIfNotExists creates a schema if it does not exist based on the given schema name.
func CreateSchemaIfNotExists(pool *pgxpool.Pool, schemaName string) (bool, error) {
	// Check if the schema exists
	var schemaExists bool
	err := pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = $1)", schemaName).Scan(&schemaExists)
	if err != nil {
		return false, err
	}

	// Check if any table exists in the public schema
	var anyTableExists bool
	err = pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' LIMIT 1)").Scan(&anyTableExists)
	if err != nil {
		return false, err
	}

	// If schema does not exist and any table exists in the public schema, create the schema
	if !schemaExists || anyTableExists {
		_, err := pool.Exec(context.Background(), fmt.Sprintf("CREATE SCHEMA %s", schemaName))
		if err != nil {
			return false, err
		}
		return true, nil
	}

	// Return false to indicate that either the schema exists or there are no tables in the public schema
	return false, nil
}

// ExecuteSQLStatement executes a SQL statement using the provided database connection pool
func ExecuteSQLStatement(pool *pgxpool.Pool, sqlStatement string) error {
	_, err := pool.Exec(context.Background(), sqlStatement)
	return err
}
