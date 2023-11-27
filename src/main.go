package main

import (
	"fmt"
	"log"
	"os"
	"schema-migration/database"
	"schema-migration/sqlstatements"
	"strings"
)

const (
	usageStatement = "usage: /usr/local/bin/schema-migration PGHOST PGDATABASE PGUSER PGPASSWORD"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatalf("Error: %v", err)
	}

}

func run(args []string) error {
	// Log the number of arguments passed
	log.Printf("Number of arguments passed: %d\n", len(args))
	// Log each argument separately
	for i, arg := range args {
		log.Printf("Argument %d: %s\n", i+1, arg)
	}
	// Assume we receive argument which contains all parameters
	argsString := strings.Join(args, " ")
	// Split the string into individual words
	argsEach := strings.Fields(argsString)
	// Check for number of arguments passed via CLI
	if len(argsEach) < 4 {
		return fmt.Errorf(usageStatement)
	}
	// Create a PostgreSQL connection configuration
	dbConfig := database.GetConfigFromArgsOrEnv(argsEach)

	// Create a PostgreSQL connection pool
	pool, err := database.NewDBPool(dbConfig)
	if err != nil {
		return fmt.Errorf("error creating connection pool: %w", err)
	}

	if os.Getenv("SCHEMA_MIGRATION") == "true" {
		// OS ENV variable related to name of the schema
		schemaName := os.Getenv(argsEach[2])

		// Check if the schema needs to be created
		if created, err := database.CreateSchemaIfNotExists(pool, schemaName); err != nil {
			return fmt.Errorf("error checking/creating schema: %w", err)
		} else if created {
			log.Printf("Creating a schema %s", schemaName)
			// List of SQL statements to execute
			sqlStatements := []string{
				sqlstatements.MoveTablesSQL,
				sqlstatements.MoveSequencesSQL,
				sqlstatements.MoveViewsSQL,
			}
			// Execute the SQL statements
			for idx, sqlStatement := range sqlStatements {
				log.Printf("Executing SQL statement %d: %s", idx+1, sqlStatement)
				if err := database.ExecuteSQLStatement(pool, sqlStatement); err != nil {
					return fmt.Errorf("error executing SQL statement %d: %w", idx+1, err)
				}
			}
			log.Println("SQL statements executed successfully.")
		} else {
			log.Printf("No execution needed since schema %s already exists", schemaName)
		}
	}
	return nil
}
