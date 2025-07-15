package tui

import (
	"fmt"
	"io"
	"os"
	"sps/internal/data"
	"strings"

	"github.com/mattn/go-runewidth"
)

// PrintTable renders data in a compact, box-style table, correctly handling CJK characters.
func PrintTable(writer io.Writer, headers []string, data [][]string) {
	if len(headers) == 0 {
		return
	}

	// 1. Calculate column widths
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = runewidth.StringWidth(h)
	}
	for _, row := range data {
		for i, cell := range row {
			width := runewidth.StringWidth(cell)
			if width > colWidths[i] {
				colWidths[i] = width
			}
		}
	}

	// 2. Create border line
	borderLine := "+"
	for _, w := range colWidths {
		borderLine += strings.Repeat("-", w+2) + "+" // +2 for padding
	}

	// 3. Print top border
	fmt.Fprintln(writer, borderLine)

	// 4. Print header
	headerLine := "|"
	for i, h := range headers {
		headerLine += " " + pad(h, colWidths[i]) + " |"
	}
	fmt.Fprintln(writer, headerLine)

	// 5. Print separator (which is the same as the border)
	fmt.Fprintln(writer, borderLine)

	// 6. Print data rows
	for _, row := range data {
		rowLine := "|"
		for i, cell := range row {
			rowLine += " " + pad(cell, colWidths[i]) + " |"
		}
		fmt.Fprintln(writer, rowLine)
	}

	// 7. Print bottom border
	fmt.Fprintln(writer, borderLine)
}

// pad pads a string to a certain visual width with spaces.
func pad(s string, width int) string {
	return s + strings.Repeat(" ", width-runewidth.StringWidth(s))
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