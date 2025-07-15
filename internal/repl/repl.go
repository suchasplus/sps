package repl

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sps/internal/adapter"
	"sps/internal/analyzer"
	"sps/internal/data"
	"sps/internal/tui"
	"strings"
)

type replContext struct {
	db        *sql.DB
	dbAdapter adapter.DBAdapter
	topN      int
	tables    []data.Table // Cache the list of tables
}

// Start initializes and runs the Read-Eval-Print Loop.
func Start(db *sql.DB, dbAdapter adapter.DBAdapter, topN int) error {
	ctx := &replContext{
		db:        db,
		dbAdapter: dbAdapter,
		topN:      topN,
	}
	// Pre-cache the table list on start
	if err := ctx.listTables(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not pre-cache table list: %v\n", err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the sps REPL. Type 'help' for commands or 'exit' to quit.")

	for {
		fmt.Print("sps> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("\nGoodbye!")
				return nil
			}
			return err
		}

		input = strings.TrimSpace(input)
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		args := parts[1:]

		if err = ctx.handleCommand(command, args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}
}

func (ctx *replContext) handleCommand(command string, args []string) error {
	switch command {
	case "help":
		printHelp()
	case "exit", "quit":
		fmt.Println("Goodbye!")
		os.Exit(0)
	case "tables", "ls":
		// Force a refresh
		return ctx.listTables()
	case "summary":
		if len(args) == 0 {
			return errors.New("usage: summary <table_name>")
		}
		return ctx.showSummary(args[0])
	case "schema":
		if len(args) == 0 {
			return errors.New("usage: schema <table_name>")
		}
		return ctx.showSchema(args[0])
	default:
		fmt.Printf("Unknown command: %s. Type 'help' for a list of commands.\n", command)
	}
	return nil
}

// findTable looks up a table by name in the cached list.
func (ctx *replContext) findTable(name string) (data.Table, bool) {
	for _, t := range ctx.tables {
		if t.Name == name {
			return t, true
		}
	}
	return data.Table{}, false
}

func (ctx *replContext) listTables() error {
	tables, err := analyzer.GetTables(ctx.db, ctx.dbAdapter)
	if err != nil {
		return err
	}
	ctx.tables = tables // Update cache

	var tableData [][]string
	headers := []string{"SCHEMA", "TABLE NAME"}
	for _, t := range tables {
		tableData = append(tableData, []string{t.Schema, t.Name})
	}
	tui.PrintTable(os.Stdout, headers, tableData)
	return nil
}

func (ctx *replContext) showSchema(tableName string) error {
	table, found := ctx.findTable(tableName)
	if !found {
		return fmt.Errorf("table %q not found", tableName)
	}
	schema, err := analyzer.GetSchema(ctx.db, ctx.dbAdapter, table)
	if err != nil {
		return err
	}
	tui.RenderSchema(schema)
	return nil
}

func (ctx *replContext) showSummary(tableName string) error {
	table, found := ctx.findTable(tableName)
	if !found {
		return fmt.Errorf("table %q not found", tableName)
	}

	totalRows, err := analyzer.GetTotalRows(ctx.db, table)
	if err != nil {
		return err
	}
	fmt.Printf("Total Rows: %d\n\n", totalRows)

	schema, err := analyzer.GetSchema(ctx.db, ctx.dbAdapter, table)
	if err != nil {
		return err
	}

	fmt.Println("[Table Schema]")
	tui.RenderSchema(schema)

	fmt.Println("\n[Column Distribution]")
	for _, col := range schema {
		dist, err := analyzer.GetColumnDistribution(ctx.db, table, col.Name, totalRows, int64(ctx.topN))
		if err != nil {
			fmt.Printf("Could not analyze column %s: %v\n", col.Name, err)
			continue
		}
		tui.RenderDistribution(dist)
	}
	return nil
}

func printHelp() {
	fmt.Println("\nAvailable commands:")
	fmt.Println("  tables, ls         - List all tables in the database.")
	fmt.Println("  summary <table_name> - Show a full summary of a table.")
	fmt.Println("  schema <table_name>  - Show the schema of a table.")
	fmt.Println("  help                 - Show this help message.")
	fmt.Println("  exit, quit           - Exit the REPL.")
}
