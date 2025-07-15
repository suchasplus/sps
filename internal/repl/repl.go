package repl

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sps/internal/adapter"
	"sps/internal/analyzer"
	"sps/internal/data"
	"sps/internal/tui"
	"strings"

	"github.com/chzyer/readline"
)

type replContext struct {
	db        *sql.DB
	dbAdapter adapter.DBAdapter
	topN      int
	tables    []data.Table // Cache the list of tables
	rl        *readline.Instance
}

// Start initializes and runs the Read-Eval-Print Loop with history and line editing.
func Start(db *sql.DB, dbAdapter adapter.DBAdapter, topN int) error {
	// Configure readline
	homeDir, _ := os.UserHomeDir()
	historyFile := filepath.Join(homeDir, ".sps_history")
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      "sps> ",
		HistoryFile: historyFile,
	})
	if err != nil {
		return err
	}
	defer rl.Close()

	ctx := &replContext{
		db:        db,
		dbAdapter: dbAdapter,
		topN:      topN,
		rl:        rl,
	}

	// Pre-cache the table list on start
	if err := ctx.listTables(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not pre-cache table list: %v\n", err)
	}
	fmt.Println("Welcome to the sps REPL. Type 'help' for commands or 'exit' to quit.")

	for {
		line, err := ctx.rl.Readline()
		if err == readline.ErrInterrupt { // Ctrl+C
			continue
		} else if err == io.EOF { // Ctrl+D
			fmt.Println("\nGoodbye!")
			break
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		ctx.rl.SaveHistory(line)

		parts := strings.Fields(line)
		command := parts[0]
		args := parts[1:]

		if err = ctx.handleCommand(command, args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}
	return nil
}

func (ctx *replContext) handleCommand(command string, args []string) error {
	switch command {
	case "help":
		printHelp()
	case "exit", "quit":
		fmt.Println("Goodbye!")
		os.Exit(0)
	case "tables", "ls":
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
	tui.PrintTable(ctx.rl.Stdout(), headers, tableData) // Use readline's writer
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