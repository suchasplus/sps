package adapter

// SQLiteAdapter implements the DBAdapter interface for SQLite.
type SQLiteAdapter struct{}

// ListTablesSQL returns the SQL query to list tables in a SQLite database.
func (a *SQLiteAdapter) ListTablesSQL() string {
	return "SELECT name FROM sqlite_master WHERE type='table';"
}

// GetTableSchemaSQL returns the SQL query to get the schema of a table in a SQLite database.
func (a *SQLiteAdapter) GetTableSchemaSQL(tableName string) string {
	return "PRAGMA table_info(" + tableName + ");"
}
