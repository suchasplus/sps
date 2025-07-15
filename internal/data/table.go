package data

// Table represents a database table with its schema and name.
type Table struct {
	Schema string
	Name   string
}

// FullyQualifiedName returns the schema-qualified name for use in queries.
func (t *Table) FullyQualifiedName() string {
	// For now, we handle the non-schema case for drivers like MySQL/SQLite.
	// A more advanced implementation would handle this based on driver type.
	if t.Schema == "" || t.Schema == "public" { // Also consider if we want to omit public
		return t.Name
	}
	return t.Schema + "." + t.Name
}
