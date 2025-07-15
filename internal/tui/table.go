package tui

import (
	"fmt"
	"io"
	"os"
	"sps/internal/data"
	"text/tabwriter"
)

// PrintTable renders data in a well-formatted table.
// headers is a slice of strings for the table header.
// data is a slice of slices of strings for the table rows.
func PrintTable(writer io.Writer, headers []string, data [][]string) {
	tw := tabwriter.NewWriter(writer, 0, 0, 3, ' ', 0)

	// Write headers
	headerLine := ""
	for i, h := range headers {
		headerLine += h
		if i < len(headers)-1 {
			headerLine += "\t"
		}
	}
	fmt.Fprintln(tw, headerLine)

	// Write data rows
	for _, row := range data {
		rowLine := ""
		for i, cell := range row {
			rowLine += cell
			if i < len(row)-1 {
				rowLine += "\t"
			}
		}
		fmt.Fprintln(tw, rowLine)
	}

	tw.Flush()
}

// RenderSchema takes a slice of ColumnSchema and prints it as a table.
func RenderSchema(schemas []data.ColumnSchema) {
	headers := []string{"COLUMN NAME", "DATA TYPE", "IS NULLABLE"}
	var data [][]string
	for _, s := range schemas {
		data = append(data, []string{s.Name, s.Type, s.IsNullable})
	}
	PrintTable(os.Stdout, headers, data)
}

// RenderDistribution takes a ColumnDistribution and prints it as a table.
func RenderDistribution(dist data.ColumnDistribution) {
	fmt.Printf("--- Column: %s ---\n", dist.ColumnName)
	headers := []string{"VALUE", "COUNT", "PERCENTAGE"}
	var tableData [][]string
	for _, v := range dist.Values {
		tableData = append(tableData, []string{
			v.Value,
			fmt.Sprintf("%d", v.Count),
			fmt.Sprintf("%.2f%%", v.Percentage),
		})
	}
	PrintTable(os.Stdout, headers, tableData)
	fmt.Println() // Add a newline for spacing
}
