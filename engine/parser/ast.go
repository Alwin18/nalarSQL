package parser

// Minimal AST shapes for MVP

type Statement interface{}

type CreateTableStmt struct {
	TableName string
	Columns   []ColumnDef
}

type ColumnDef struct {
	Name string
	Type string
}

type InsertStmt struct {
	Table   string
	Columns []string
	Values  []any
}

type SelectStmt struct {
	Table   string
	Columns []string // empty or ["*"] = all
}
