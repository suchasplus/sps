package analyzer

import (
	"database/sql"
	"fmt"
	"sps/internal/data"
)

// GetTotalRows counts the total number of rows in a table.
func GetTotalRows(db *sql.DB, table data.Table) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table.FullyQualifiedName())
	err := db.QueryRow(query).Scan(&count)
	return count, err
}

// GetColumnDistribution calculates the value distribution for a single column.
func GetColumnDistribution(db *sql.DB, table data.Table, columnName string, totalRows, topN int64) (data.ColumnDistribution, error) {
	dist := data.ColumnDistribution{ColumnName: columnName}

	query := fmt.Sprintf(
		"SELECT %q, COUNT(*) as count FROM %s GROUP BY %q ORDER BY count DESC LIMIT %d",
		columnName, table.FullyQualifiedName(), columnName, topN,
	)

	rows, err := db.Query(query)
	if err != nil {
		return dist, err
	}
	defer rows.Close()

	for rows.Next() {
		var value sql.NullString // Use NullString to handle potential NULL values
		var count int64
		if err := rows.Scan(&value, &count); err != nil {
			return dist, err
		}

		valStr := "NULL"
		if value.Valid {
			valStr = value.String
		}

		percentage := 0.0
		if totalRows > 0 {
			percentage = (float64(count) / float64(totalRows)) * 100
		}

		dist.Values = append(dist.Values, data.ValueCount{
			Value:      valStr,
			Count:      count,
			Percentage: percentage,
		})
	}

	return dist, nil
}