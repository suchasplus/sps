package adapter

// MySQLAdapter implements the DBAdapter interface for MySQL.
type MySQLAdapter struct{}

// ListTablesSQL returns the SQL query to list tables in a MySQL database.
func (a *MySQLAdapter) ListTablesSQL() string {
	return "SHOW TABLES;"
}

// GetTableSchemaSQL returns the SQL query to get the schema of a table in a MySQL database.
func (a *MySQLAdapter) GetTableSchemaSQL(tableName string) string {
	// Using DESCRIBE is simpler than querying INFORMATION_SCHEMA in MySQL
	return "DESCRIBE " + tableName + ";"
}
