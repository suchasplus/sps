# sps - Database Summary CLI - Design Document (v2)

## 1. 简介 (Introduction)

本项目旨在创建一个基于Go语言的命令行工具（CLI），用于快速分析和汇总数据库中的表数据。其灵感来源于R语言中强大的`summary()`函数，旨在为开发人员和数据分析师提供一个无需离开终端即可快速洞察数据概况的工具。

用户可以通过直接传入DSN字符串或从文件中安全加载DSN来连接数据库。连接后，工具支持两种模式：直接对指定表进行分析，或进入一个交互式的REPL（Read-Eval-Print Loop）环境，允许用户连续执行多个命令（如列出表、查看表结构、分析数据分布）进行探索性分析。

## 2. 核心功能 (Core Features)

- **多数据库支持**: 通过标准的DSN支持连接到以下主流数据库：
    - PostgreSQL
    - MySQL
    - SQLite
    - Microsoft SQL Server
- **安全灵活的DSN处理**:
    - 支持将DSN作为直接参数传入。
    - 支持通过 `-f, --file` 标志从文件中读取DSN，避免密码等敏感信息在命令行历史中泄露。
- **双模式操作**:
    - **直接执行模式**: 当用户提供DSN并用 `-t, --table` 指定表名时，程序直接输出分析结果并退出。
    - **交互式REPL模式**: 当用户仅提供DSN时，程序进入一个REPL环境，支持连续的数据库探索。
- **表结构摘要 (Schema Summary)**:
    - 显示表的完整结构信息，包括列名、数据类型、是否可为空等。
- **数据分布分析 (Data Distribution Analysis)**:
    - 计算并展示表中总行数。
    - 对每一列计算唯一值的频率和占比。
- **智能输出控制**:
    - **简洁输出**: 默认情况下，连接成功等非关键信息不会显示。使用 `-v, --verbose` 标志可启用详细输出。
    - **高基数列处理**: 对于唯一值过多的列，默认仅显示频率最高的前N个值。

## 3. 架构设计 (Architecture)

- **CLI框架 (CLI Framework)**:
    - 使用 `github.com/urfave/cli/v3` 构建CLI。它是一个强大且富有表现力的框架，非常适合构建结构清晰的Go命令行应用。
- **数据库抽象层 (Database Abstraction)**:
    - 核心使用Go内置的 `database/sql` 包和相应的数据库驱动。
    - 设计一个`DBAdapter`接口来封装不同数据库的SQL方言差异（如获取表列表的查询）。
- **核心分析逻辑 (Core Analysis Logic)**:
    - 独立于CLI和数据库层，负责执行SQL查询并计算统计数据。
- **REPL引擎 (REPL Engine)**:
    - 一个循环，使用Go标准库（如 `bufio.Scanner`）读取用户输入。
    - 对输入进行简单的解析，匹配到预定义的REPL命令（如 `summary`, `tables`）并调用相应的分析函数。
- **用户交互 (User Interaction)**:
    - 在直接执行模式下，如果未指定表名，将使用 `github.com/manifoldco/promptui` 提供交互式表选择。
- **输出格式化 (Output Formatting)**:
    - 使用 `text/tabwriter` 来确保输出的表格对齐整洁。

## 4. 实现细节 (Implementation Details)

### 4.1. DSN处理逻辑

1.  程序启动时，首先检查是否存在 `-f, --file` 标志。
2.  如果存在，从指定文件中读取DSN内容。
3.  如果不存在，则从命令行的第一个参数中获取DSN。
4.  如果两种方式都无法获取DSN，则报错退出。

### 4.2. 数据库方言差异

- **列出所有表 (List Tables)**:
    - **PostgreSQL**: `SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname NOT IN ('pg_catalog', 'information_schema');`
    - **MySQL**: `SHOW TABLES;`
    - **SQLite**: `SELECT name FROM sqlite_master WHERE type='table';`
    - **MSSQL**: `SELECT table_name FROM information_schema.tables WHERE table_type = 'BASE TABLE';`

### 4.3. 命令行接口 (CLI Usage)

```bash
# Usage
sps [command options] [DSN]

# --- Flags ---
# DSN Input (choose one)
DSN Argument: "user:pass@tcp(host:port)/db"
-f, --file FILE: Load DSN from FILE

# Mode Control
-t, --table TABLE: Specify table and run in direct mode

# Output & Analysis Control
-v, --verbose: Enable verbose logging (e.g., show connection success message)
-l, --limit N:   Sample N rows for analysis (default: all)
-n, --top N:     Show top N frequent values for high-cardinality columns (default: 20)

# --- Examples ---

# 1. Enter REPL mode using DSN from a file
sps -f /path/to/my.dsn

# 2. Enter REPL mode using DSN as an argument
sps "user:password@tcp(127.0.0.1:3306)/mydatabase"

# 3. Direct analysis on a specific table (quiet by default)
sps -t employees "user:pass@tcp(host:port)/employees"
```

## 5. REPL模式详解 (REPL Mode In-Depth)

当仅提供DSN时，用户会进入一个以数据库名作为提示符的交互环境。

### REPL命令

- `summary <table_name>`: 对指定表进行完整的结构和数据分布分析。
- `schema <table_name>`: 仅显示指定表的结构。
- `tables` or `ls`: 列出当前数据库中的所有表。
- `help`: 显示所有可用的REPL命令。
- `exit` or `quit`: 退出REPL会话。

### REPL会话示例

```
$ sps -f ./my.dsn -v
✅ Connected to MySQL database 'employees' successfully.

employees> ls
+-----------+
| TABLES    |
+-----------+
| employees |
| departments|
| salaries  |
+-----------+

employees> summary salaries
Table Summary: `salaries`
Total Rows: 2,844,047

[Table Schema]
+-----------+---------+------------+
| COLUMN    | TYPE    | NULLABLE   |
+-----------+---------+------------+
| emp_no    | int     | NO         |
| salary    | int     | NO         |
| from_date | date    | NO         |
| to_date   | date    | NO         |
+-----------+---------+------------+

[Column Distribution]
... (analysis output) ...

employees> quit
Goodbye!
```

## 6. 未来可能的增强 (Future Enhancements)

- **更多统计指标**: 对于数值型列，增加min, max, mean, stddev, and quartiles等统计摘要。
- **数据可视化**: 为数值列生成小型的文本直方图（histogram）。
- **更多输出格式**: 支持将结果导出为JSON或CSV格式，方便与其他工具集成。
- **配置文件**: 支持通过`~/.config/sps.toml`等文件预设DSN别名，方便快速连接。