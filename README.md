# sps - SQL Profile & Summary 123

`sps` is a command-line tool inspired by R's `summary()` function, designed for developers and data analysts to quickly profile and summarize database tables directly from the terminal.

![sps-demo](https://user-images.githubusercontent.com/12345/some-image.gif) <!-- Placeholder for a demo GIF -->

## Features

-   **Multi-Database Support**: Natively connects to PostgreSQL, MySQL, and SQLite.
-   **Secure & Flexible Connection**: Provide your database DSN as a direct argument or safely from a file using the `-f` flag to avoid exposing credentials in your shell history.
-   **Dual-Mode Operation**:
    -   **Direct Mode**: Get a quick, one-off summary of a single table using the `-t` flag.
    -   **REPL Mode**: An interactive shell for exploring your database, complete with command history (up/down arrows) and line editing.
-   **Comprehensive Summary**:
    -   **Table Schema**: Displays column names, data types, and nullability.
    -   **Row Count**: Shows the total number of rows in the table.
    -   **Data Distribution**: For each column, it calculates the frequency and percentage of unique values.
-   **User-Friendly Output**: Renders data in clean, box-style tables that correctly align CJK and other wide characters.

## Installation

### From Source

Ensure you have Go installed on your system.

```bash
git clone https://github.com/suchasplus/sps.git
cd sps
make build
# The 'sps' binary will be in the root directory.
# You can move it to a directory in your PATH, e.g., /usr/local/bin
```

### Using `go install`

You can also install the binary directly using `go install`:
```bash
go install github.com/suchasplus/sps@latest
```

## Usage

### Direct Analysis Mode

To get a quick summary of a specific table, use the `-t` flag.

```bash
# Using DSN as an argument
sps -t employees "user:pass@tcp(host:port)/employees"

# Using DSN from a file for better security
sps -f /path/to/my.dsn -t salaries
```

### REPL Mode

For interactive exploration, simply provide the DSN without the `-t` flag. This will launch the REPL.

```bash
sps "postgres://user:pass@host:5432/mydatabase?sslmode=disable"
```

Once inside the REPL (`sps>`), you can use the following commands:

-   `ls` or `tables`: List all schemas and tables.
-   `schema <table_name>`: Show the schema for a specific table.
-   `summary <table_name>`: Provide a full summary for a specific table.
-   `help`: Show a list of available commands.
-   `exit` or `quit`: Exit the REPL session.

The REPL supports command history, which is saved to `~/.sps_history`.

## Development

This project uses a `Makefile` to streamline common development tasks.

-   `make build`: Build the binary.
-   `make test`: Run all unit tests.
-   `make lint`: Lint the codebase using `golangci-lint`.
-   `make clean`: Clean build artifacts.
-   `make help`: See all available commands.
