package analyzer

import (
	"database/sql"
	"regexp"
	"sps/internal/data"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// Helper function to create a mock DB and adapter for tests
func newMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestGetTotalRows(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	tableName := "users"
	expectedRows := int64(150)

	rows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(expectedRows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM " + tableName)).WillReturnRows(rows)

	count, err := GetTotalRows(db, tableName)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if count != expectedRows {
		t.Errorf("expected row count %d, got %d", expectedRows, count)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetColumnDistribution(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	tableName := "users"
	columnName := "country"
	totalRows := int64(10)

	mockRows := sqlmock.NewRows([]string{columnName, "count"}).
		AddRow("USA", 5).
		AddRow("Canada", 3).
		AddRow("Mexico", 2)

	query := "SELECT country, COUNT(*) as count FROM users GROUP BY country ORDER BY count DESC LIMIT 10"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(mockRows)

	dist, err := GetColumnDistribution(db, tableName, columnName, totalRows, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dist.Values) != 3 {
		t.Fatalf("expected 3 values, got %d", len(dist.Values))
	}

	if dist.Values[0].Value != "USA" || dist.Values[0].Count != 5 || dist.Values[0].Percentage != 50.0 {
		t.Errorf("unexpected distribution for USA: %+v", dist.Values[0])
	}
	if dist.Values[1].Value != "Canada" || dist.Values[1].Count != 3 || dist.Values[1].Percentage != 30.0 {
		t.Errorf("unexpected distribution for Canada: %+v", dist.Values[1])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetSchema(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	// Mock adapter that returns a generic schema query
	mockAdapter := &mockAdapter{schemaSQL: "DESCRIBE users"}
	tableName := "users"

	mockRows := sqlmock.NewRows([]string{"Field", "Type", "Null"}).
		AddRow("id", "int", "NO").
		AddRow("name", "varchar(255)", "YES")

	mock.ExpectQuery(regexp.QuoteMeta("DESCRIBE users")).WillReturnRows(mockRows)

	schema, err := GetSchema(db, mockAdapter, tableName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedSchema := []data.ColumnSchema{
		{Name: "id", Type: "int", IsNullable: "NO"},
		{Name: "name", Type: "varchar(255)", IsNullable: "YES"},
	}

	if len(schema) != len(expectedSchema) {
		t.Fatalf("expected schema length %d, got %d", len(expectedSchema), len(schema))
	}

	for i, col := range schema {
		if col.Name != expectedSchema[i].Name || col.Type != expectedSchema[i].Type || col.IsNullable != expectedSchema[i].IsNullable {
			t.Errorf("expected schema %+v, got %+v", expectedSchema[i], col)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// mockAdapter is a simple adapter for testing purposes.
type mockAdapter struct {
	listSQL   string
	schemaSQL string
}

func (m *mockAdapter) ListTablesSQL() string {
	return m.listSQL
}

func (m *mockAdapter) GetTableSchemaSQL(tableName string) string {
	return m.schemaSQL
}
