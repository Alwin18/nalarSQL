package planner

import (
	"github.com/Alwin18/nalarSQL/engine/parser"
	"github.com/Alwin18/nalarSQL/engine/storage"
)

// Planner converts AST -> Plan (very small abstraction)

type Planner struct {
	store *storage.Store
}

func NewPlanner(s *storage.Store) *Planner { return &Planner{store: s} }

type Plan interface{}

type PlanCreateTable struct {
	Stmt *parser.CreateTableStmt
}

type PlanInsert struct {
	Stmt *parser.InsertStmt
}

type PlanSelect struct {
	Stmt *parser.SelectStmt
}

func (p *Planner) Plan(stmt parser.Statement) (Plan, error) {
	switch s := stmt.(type) {
	case *parser.CreateTableStmt:
		return &PlanCreateTable{Stmt: s}, nil
	case *parser.InsertStmt:
		return &PlanInsert{Stmt: s}, nil
	case *parser.SelectStmt:
		return &PlanSelect{Stmt: s}, nil
	default:
		return nil, ErrUnsupportedPlan
	}
}
