package parser

// Minimal AST shapes for MVP

type Statement interface {
	stmt() // marker method
}

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

type UpdateStmt struct {
	Table       string
	Set         map[string]any
	WhereColumn string
	WhereValue  any
}

type DeleteStmt struct {
	Table       string
	WhereColumn string
	WhereValue  any
}

// Implement Statement interface marker methods
func (*CreateTableStmt) stmt() {}
func (*InsertStmt) stmt()      {}
func (*SelectStmt) stmt()      {}
func (*UpdateStmt) stmt()      {}
func (*DeleteStmt) stmt()      {}
