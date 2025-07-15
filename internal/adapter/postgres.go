package adapter

// PostgresAdapter implements the DBAdapter interface for PostgreSQL.
type PostgresAdapter struct{}

// ListTablesSQL returns the SQL query to list tables and their schemas in a PostgreSQL database.
func (a *PostgresAdapter) ListTablesSQL() string {
	return "SELECT schemaname, tablename FROM pg_catalog.pg_tables WHERE schemaname NOT IN ('pg_catalog', 'information_schema');"
}

// GetTableSchemaSQL returns the SQL query to get the schema of a table in a PostgreSQL database.
func (a *PostgresAdapter) GetTableSchemaSQL(tableName string) string {
	return "SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = '" + tableName + "';"
}
