# nalarSQL - Interactive SQL Database

A simple SQL database engine with interactive CLI written in Go.

## Features

- ✅ **CREATE TABLE** - Create tables with column definitions
- ✅ **INSERT INTO** - Insert data into tables
- ✅ **SELECT** - Query data with column projection and table display
- ✅ **UPDATE** - Update records with WHERE clause
- ✅ **DELETE** - Delete records with WHERE clause
- ✅ **Interactive CLI** - REPL interface for running SQL commands
- ✅ **Beautiful Output** - Color-coded table display with proper formatting

## Building

```bash
go build -o nalarSql .
```

## Running

```bash
./nalarSql
```

Or run the demo:
```bash
chmod +x demo.sh
./demo.sh
```

## Usage Examples

### Starting the CLI

```
$ ./nalarSql
┌─────────────────────────────────────┐
│   Welcome to nalarSQL Database!    │
│   Type 'exit' or 'quit' to exit    │
└─────────────────────────────────────┘

nalarSQL>
```

### Creating a Table

```sql
nalarSQL> CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, age INTEGER);
✅ Query executed successfully
```

### Inserting Data

```sql
nalarSQL> INSERT INTO users (id, name, age) VALUES (1, 'Alice', 30);
✅ Row inserted with ID: 1

nalarSQL> INSERT INTO users (id, name, age) VALUES (2, 'Bob', 25);
✅ Row inserted with ID: 2
```

### Selecting Data (with Beautiful Table Display!)

```sql
nalarSQL> SELECT * FROM users;
┌─────┬───────┬────┐
│ age │ name  │ id │
├─────┼───────┼────┤
│ 30  │ Alice │ 1  │
│ 25  │ Bob   │ 2  │
└─────┴───────┴────┘
2 rows returned
```

### Updating Data

```sql
nalarSQL> UPDATE users SET age = 31 WHERE id = 1;
✅ 1 row updated

nalarSQL> SELECT * FROM users;
┌─────┬───────┬────┐
│ age │ name  │ id │
├─────┼───────┼────┤
│ 31  │ Alice │ 1  │
│ 25  │ Bob   │ 2  │
└─────┴───────┴────┘
2 rows returned
```

### Deleting Data

```sql
nalarSQL> DELETE FROM users WHERE id = 2;
✅ 1 row deleted

nalarSQL> SELECT * FROM users;
┌─────┬───────┬────┐
│ age │ name  │ id │
├─────┼───────┼────┤
│ 31  │ Alice │ 1  │
└─────┴───────┴────┘
1 row returned
```

### Exiting

```sql
nalarSQL> exit
👋 Goodbye!
```

## Architecture

```
nalarSQL/
├── main.go              # Interactive CLI entry point
├── engine/
│   ├── engine.go        # Main engine facade
│   ├── parser/          # SQL parser & lexer
│   │   ├── lexer.go    # Tokenizer
│   │   ├── parser.go   # SQL parser
│   │   └── ast.go      # AST definitions
│   ├── planner/         # Query planner
│   │   └── planner.go
│   ├── executor/        # Query executor
│   │   └── executor.go
│   └── storage/         # Storage engine
│       └── store.go     # File-based storage
└── .data/               # Database files (auto-created)
```

## Storage Format

Tables are stored as JSON files in the `.data/` directory:
- First line: Schema header with column definitions
- Following lines: One JSON object per row

## Supported SQL

### CREATE TABLE
```sql
CREATE TABLE table_name (
    column1 TYPE [constraints],
    column2 TYPE [constraints],
    ...
);
```

Supported types: `INTEGER`, `TEXT`
Constraints are parsed but not enforced (for compatibility)

### INSERT
```sql
INSERT INTO table_name (col1, col2, ...) VALUES (val1, val2, ...);
```

### SELECT
```sql
SELECT * FROM table_name;
SELECT col1, col2 FROM table_name;
```

### UPDATE
```sql
UPDATE table_name SET col1 = val1, col2 = val2 WHERE column = value;
```

### DELETE
```sql
DELETE FROM table_name WHERE column = value;
```

## Limitations

- No JOIN support
- WHERE clause only supports simple equality (col = value)
- No ORDER BY, GROUP BY, LIMIT
- No transactions
- Single-threaded
- File locking is basic (not suitable for concurrent access)

## License

MIT