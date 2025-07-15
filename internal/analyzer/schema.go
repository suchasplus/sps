package analyzer

import (
	"database/sql"
	"fmt"
	"sps/internal/adapter"
	"sps/internal/data"
)

// GetSchema retrieves the schema for a given table.
// Note: The parsing logic here is simplified and might need to be adjusted
// based on the exact output of each database's schema query.
func GetSchema(db *sql.DB, dbAdapter adapter.DBAdapter, table data.Table) ([]data.ColumnSchema, error) {
	rows, err := db.Query(dbAdapter.GetTableSchemaSQL(table.Name))
	if err != nil {
		return nil, fmt.Errorf("failed to query schema for table %s: %w", table.Name, err)
	}
	defer rows.Close()

	var columns []data.ColumnSchema
	
	// This is a generic approach; specific adapters might need custom parsing.
	// For now, we assume a common structure for simplicity.
	// MySQL: Field, Type, Null, Key, Default, Extra
	// PostgreSQL: column_name, data_type, is_nullable
	// SQLite: cid, name, type, notnull, dflt_value, pk
	
	// A more robust solution would involve parsing based on the driver type.
	// We will start with a simple version.
	
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		// Create a slice of interface{}'s to represent a row, and a slice of pointers to each
		vals := make([]interface{}, len(cols))
		valPtrs := make([]interface{}, len(cols))
		for i := 0; i < len(cols); i++ {
			valPtrs[i] = &vals[i]
		}

		if err := rows.Scan(valPtrs...); err != nil {
			return nil, err
		}

		var col data.ColumnSchema
		// This mapping is a simplification.
		// It will be improved later based on driver type.
		col.Name = fmt.Sprintf("%s", vals[0])
		col.Type = fmt.Sprintf("%s", vals[1])
		if len(vals) > 2 {
			col.IsNullable = fmt.Sprintf("%s", vals[2])
		}
		columns = append(columns, col)
	}

	return columns, nil
}
