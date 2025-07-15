package analyzer

import (
	"database/sql"
	"fmt"
	"sps/internal/adapter"
	"sps/internal/data"
	"strings"
)

// GetTables retrieves a list of tables from the connected database.
func GetTables(db *sql.DB, dbAdapter adapter.DBAdapter) ([]data.Table, error) {
	rows, err := db.Query(dbAdapter.ListTablesSQL())
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var tables []data.Table
	for rows.Next() {
		// This logic handles both single-column (MySQL/SQLite) and two-column (Postgres) results.
		if len(cols) == 1 {
			var name string
			if err := rows.Scan(&name); err != nil {
				return nil, fmt.Errorf("failed to scan table name: %w", err)
			}
			// For drivers without schemas in this query, we can leave it blank.
			tables = append(tables, data.Table{Name: name})
		} else if len(cols) >= 2 && (strings.ToLower(cols[0]) == "schemaname" || strings.ToLower(cols[1]) == "tablename") {
			var schema, name string
			if err := rows.Scan(&schema, &name); err != nil {
				return nil, fmt.Errorf("failed to scan schema and table name: %w", err)
			}
			tables = append(tables, data.Table{Schema: schema, Name: name})
		}
	}

	return tables, nil
}
