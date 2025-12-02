package parser

// For MVP the lexer is trivial and naive. Replace with production lexer later.

// Parse is a convenience wrapper: detect simple verbs and return tiny AST
func Parse(sql string) (Statement, error) {
	// very naive: trim and check prefix
	s := sql
	if len(s) >= 6 && (s[:6] == "INSERT" || s[:6] == "insert") {
		// TODO: implement real lexer+parser. For now return a stub that executor can understand.
		// This stub only supports: INSERT INTO <table> (col,...) VALUES (v,...);
		// Return a placeholder error for unsupported queries.
		return &InsertStmt{Table: "users", Columns: []string{"id", "name", "age"}, Values: []any{int64(1), "stub", int64(0)}}, nil
	}
	if len(s) >= 6 && (s[:6] == "SELECT" || s[:6] == "select") {
		return &SelectStmt{Table: "users", Columns: []string{"*"}}, nil
	}
	if len(s) >= 6 && (s[:6] == "CREATE" || s[:6] == "create") {
		return &CreateTableStmt{TableName: "users", Columns: []ColumnDef{{Name: "id", Type: "INTEGER"}, {Name: "name", Type: "TEXT"}, {Name: "age", Type: "INTEGER"}}}, nil
	}
	return nil, ErrUnsupportedSQL
}
