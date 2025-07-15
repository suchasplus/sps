package adapter

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	// _ "github.com/denisenkom/go-mssqldb" // MSSQL driver can be added here
)

type DBAdapter interface {
	ListTablesSQL() string
	GetTableSchemaSQL(tableName string) string
	// Other database-specific methods can be added here
}

// DetectDriver determines the database driver name from a DSN string.
func DetectDriver(dsn string) (string, error) {
	switch {
	case strings.HasPrefix(dsn, "postgres://"), strings.HasPrefix(dsn, "postgresql://"):
		return "postgres", nil
	case strings.Contains(dsn, "@tcp("):
		return "mysql", nil
	case strings.HasSuffix(dsn, ".db"), strings.HasSuffix(dsn, ".sqlite"), strings.HasSuffix(dsn, ".sqlite3"):
		return "sqlite3", nil
	// Add case for mssql if needed
	default:
		return "", fmt.Errorf("could not detect driver for DSN")
	}
}

// Connect establishes a database connection using the appropriate driver.
func Connect(dsn string) (*sql.DB, string, error) {
	driverName, err := DetectDriver(dsn)
	if err != nil {
		return nil, "", err
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, "", err
	}

	if err = db.Ping(); err != nil {
		return nil, "", fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, driverName, nil
}

// NewAdapter returns a DBAdapter for the given driver name.
func NewAdapter(driverName string) (DBAdapter, error) {
	switch driverName {
	case "mysql":
		return &MySQLAdapter{}, nil
	case "postgres":
		return &PostgresAdapter{}, nil
	case "sqlite3":
		return &SQLiteAdapter{}, nil
	default:
		return nil, fmt.Errorf("unsupported driver: %s", driverName)
	}
}
