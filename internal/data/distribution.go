package data

// ValueCount holds the frequency and percentage of a specific value.
type ValueCount struct {
	Value      string
	Count      int64
	Percentage float64
}

// ColumnDistribution holds the distribution analysis for a single column.
type ColumnDistribution struct {
	ColumnName string
	Values     []ValueCount
}

// TableAnalysis holds the complete analysis for a table.
type TableAnalysis struct {
	TableName    string
	TotalRows    int64
	Schema       []ColumnSchema
	Distribution []ColumnDistribution
}
