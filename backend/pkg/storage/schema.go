package storage

import (
	"context"
	_ "embed"
	"fmt"
)

//go:embed sql/create-table-tests.sql
var schemaTests string

//go:embed sql/create-table-summary.sql
var schemaSummary string

//go:embed sql/create-table-suite-summary.sql
var schemaSuiteSummary string

//go:embed sql/create-table-metadata.sql
var schemaMetadata string

func CreateSchema(ctx context.Context) error {
	conn := db()
	if conn == nil {
		return fmt.Errorf("failed to read db")
	}
	defer conn.Close()

	_, err := conn.QueryxContext(ctx, schemaTests)
	if err != nil {
		return err
	}

	_, err = conn.QueryxContext(ctx, schemaSummary)
	if err != nil {
		return err
	}

	_, err = conn.QueryxContext(ctx, schemaSuiteSummary)
	if err != nil {
		return err
	}

	_, err = conn.QueryxContext(ctx, schemaMetadata)
	return err
}
