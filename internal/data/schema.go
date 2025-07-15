package data

// ColumnSchema holds the schema information for a single table column.
type ColumnSchema struct {
	Name       string
	Type       string
	IsNullable string // Using string to accommodate different DB responses ("YES", "NO", "0", "1")
}
