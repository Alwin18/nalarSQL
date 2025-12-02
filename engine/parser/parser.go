package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	l     *Lexer
	cur   Token
	peekT Token
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	p.next()
	p.next()
	return p
}

func (p *Parser) next() {
	p.cur = p.peekT
	p.peekT = p.l.NextToken()
}

func (p *Parser) expect(tt TokenType, val string) error {
	if p.cur.Type != tt || (val != "" && strings.ToUpper(p.cur.Value) != val) {
		return fmt.Errorf("expected %s '%s', got %v '%s'", tt, val, p.cur.Type, p.cur.Value)
	}
	p.next()
	return nil
}

// Parse only INSERT for now
func (p *Parser) ParseStatement() (Statement, error) {
	// move to first token
	if p.cur.Type == TokEOF {
		return nil, ErrUnsupportedSQL
	}
	if p.cur.Type == TokKeyword {
		switch p.cur.Value {
		case "INSERT":
			return p.parseInsert()
		case "SELECT":
			return p.parseSelect()
		case "UPDATE":
			return p.parseUpdate()
		case "DELETE":
			return p.parseDelete()
		case "CREATE":
			// delegate to previous simplistic behavior
			return p.parseCreate()
		}
	}
	return nil, ErrUnsupportedSQL
}

func (p *Parser) parseCreate() (*CreateTableStmt, error) {
	// very simple: CREATE TABLE name (col TYPE [constraints], ...)
	if err := p.expect(TokKeyword, "CREATE"); err != nil {
		return nil, err
	}
	if err := p.expect(TokKeyword, "TABLE"); err != nil {
		return nil, err
	}
	if p.cur.Type != TokIdent {
		return nil, fmt.Errorf("expected table name")
	}
	name := p.cur.Value
	p.next()
	if err := p.expect(TokLParen, ""); err != nil {
		return nil, err
	}
	cols := []ColumnDef{}
	for {
		if p.cur.Type != TokIdent {
			return nil, fmt.Errorf("expected column name")
		}
		col := p.cur.Value
		p.next()
		if p.cur.Type != TokIdent {
			return nil, fmt.Errorf("expected column type")
		}
		typ := strings.ToUpper(p.cur.Value)
		p.next()

		// Skip constraints like PRIMARY KEY, NOT NULL, etc.
		for p.cur.Type == TokIdent || p.cur.Type == TokKeyword {
			constraintUpper := strings.ToUpper(p.cur.Value)
			if constraintUpper == "PRIMARY" || constraintUpper == "NOT" ||
				constraintUpper == "NULL" || constraintUpper == "KEY" ||
				constraintUpper == "UNIQUE" || constraintUpper == "AUTO_INCREMENT" ||
				constraintUpper == "AUTOINCREMENT" || constraintUpper == "DEFAULT" {
				p.next()
				// For DEFAULT, skip the value as well
				if constraintUpper == "DEFAULT" {
					p.next()
				}
			} else {
				break
			}
		}

		cols = append(cols, ColumnDef{Name: col, Type: typ})
		if p.cur.Type == TokRParen {
			p.next()
			break
		}
		if err := p.expect(TokComma, ""); err != nil {
			return nil, err
		}
	}
	return &CreateTableStmt{TableName: name, Columns: cols}, nil
}

func (p *Parser) parseInsert() (*InsertStmt, error) {
	if err := p.expect(TokKeyword, "INSERT"); err != nil {
		return nil, err
	}
	if err := p.expect(TokKeyword, "INTO"); err != nil {
		return nil, err
	}

	if p.cur.Type != TokIdent {
		return nil, fmt.Errorf("expected table name, got %v", p.cur)
	}
	table := p.cur.Value
	p.next()

	if err := p.expect(TokLParen, ""); err != nil {
		return nil, err
	}
	cols := []string{}
	for {
		if p.cur.Type != TokIdent {
			return nil, fmt.Errorf("expected column name")
		}
		cols = append(cols, p.cur.Value)
		p.next()
		if p.cur.Type == TokRParen {
			p.next()
			break
		}
		if err := p.expect(TokComma, ""); err != nil {
			return nil, err
		}
	}

	if err := p.expect(TokKeyword, "VALUES"); err != nil {
		return nil, err
	}
	if err := p.expect(TokLParen, ""); err != nil {
		return nil, err
	}

	vals := []any{}
	for {
		switch p.cur.Type {
		case TokNumber:
			n, _ := strconv.ParseInt(p.cur.Value, 10, 64)
			vals = append(vals, n)
		case TokString:
			vals = append(vals, p.cur.Value)
		default:
			return nil, fmt.Errorf("unexpected value token %v", p.cur)
		}
		p.next()
		if p.cur.Type == TokRParen {
			p.next()
			break
		}
		if err := p.expect(TokComma, ""); err != nil {
			return nil, err
		}
	}

	return &InsertStmt{Table: table, Columns: cols, Values: vals}, nil
}

func (p *Parser) parseSelect() (*SelectStmt, error) {
	if err := p.expect(TokKeyword, "SELECT"); err != nil {
		return nil, err
	}
	cols := []string{}
	if p.cur.Type == TokStar {
		cols = append(cols, "*")
		p.next()
	} else {
		for {
			if p.cur.Type != TokIdent {
				return nil, fmt.Errorf("expected column in select")
			}
			cols = append(cols, p.cur.Value)
			p.next()
			if p.cur.Type == TokComma {
				p.next()
				continue
			}
			break
		}
	}
	if err := p.expect(TokKeyword, "FROM"); err != nil {
		return nil, err
	}
	if p.cur.Type != TokIdent {
		return nil, fmt.Errorf("expected table name in select")
	}
	table := p.cur.Value
	p.next()
	return &SelectStmt{Table: table, Columns: cols}, nil
}

func (p *Parser) parseUpdate() (*UpdateStmt, error) {
	// UPDATE <table> SET col = val [, ...] [WHERE col = val]
	if err := p.expect(TokKeyword, "UPDATE"); err != nil {
		return nil, err
	}
	if p.cur.Type != TokIdent {
		return nil, fmt.Errorf("expected table name in update")
	}
	table := p.cur.Value
	p.next()
	if err := p.expect(TokKeyword, "SET"); err != nil {
		return nil, err
	}
	set := map[string]any{}
	for {
		if p.cur.Type != TokIdent {
			return nil, fmt.Errorf("expected column in set")
		}
		col := p.cur.Value
		p.next()
		if err := p.expect(TokEqual, ""); err != nil {
			return nil, err
		}
		switch p.cur.Type {
		case TokNumber:
			n, _ := strconv.ParseInt(p.cur.Value, 10, 64)
			set[col] = n
		case TokString:
			set[col] = p.cur.Value
		default:
			return nil, fmt.Errorf("unexpected token in set %v", p.cur)
		}
		p.next()
		if p.cur.Type == TokComma {
			p.next()
			continue
		}
		break
	}
	// optional WHERE
	whereCol := ""
	var whereVal any = nil
	if p.cur.Type == TokKeyword && p.cur.Value == "WHERE" {
		p.next()
		if p.cur.Type != TokIdent {
			return nil, fmt.Errorf("expected column in where")
		}
		whereCol = p.cur.Value
		p.next()
		if err := p.expect(TokEqual, ""); err != nil {
			return nil, err
		}
		switch p.cur.Type {
		case TokNumber:
			n, _ := strconv.ParseInt(p.cur.Value, 10, 64)
			whereVal = n
		case TokString:
			whereVal = p.cur.Value
		default:
			return nil, fmt.Errorf("unexpected token in where %v", p.cur)
		}
		p.next()
	}
	return &UpdateStmt{Table: table, Set: set, WhereColumn: whereCol, WhereValue: whereVal}, nil
}

func (p *Parser) parseDelete() (*DeleteStmt, error) {
	// DELETE FROM <table> [WHERE col = val]
	if err := p.expect(TokKeyword, "DELETE"); err != nil {
		return nil, err
	}
	if err := p.expect(TokKeyword, "FROM"); err != nil {
		return nil, err
	}
	if p.cur.Type != TokIdent {
		return nil, fmt.Errorf("expected table name in delete")
	}
	table := p.cur.Value
	p.next()
	whereCol := ""
	var whereVal any = nil
	if p.cur.Type == TokKeyword && p.cur.Value == "WHERE" {
		p.next()
		if p.cur.Type != TokIdent {
			return nil, fmt.Errorf("expected column in where")
		}
		whereCol = p.cur.Value
		p.next()
		if err := p.expect(TokEqual, ""); err != nil {
			return nil, err
		}
		switch p.cur.Type {
		case TokNumber:
			n, _ := strconv.ParseInt(p.cur.Value, 10, 64)
			whereVal = n
		case TokString:
			whereVal = p.cur.Value
		default:
			return nil, fmt.Errorf("unexpected token in where %v", p.cur)
		}
		p.next()
	}
	return &DeleteStmt{Table: table, WhereColumn: whereCol, WhereValue: whereVal}, nil
}
