package engine

import (
	"path/filepath"

	"github.com/Alwin18/nalarSQL/engine/executor"
	"github.com/Alwin18/nalarSQL/engine/parser"
	"github.com/Alwin18/nalarSQL/engine/planner"
	"github.com/Alwin18/nalarSQL/engine/storage"
)

// Engine is the facade that ties parser, planner, storage and executor
type Engine struct {
	stor *storage.Store
	pl   *planner.Planner
	ex   *executor.Executor
}

// NewEngine opens/creates data dir
func NewEngine(dataDir string) (*Engine, error) {
	st, err := storage.NewStore(filepath.Clean(dataDir))
	if err != nil {
		return nil, err
	}
	pl := planner.NewPlanner(st)
	exec := executor.NewExecutor(st)
	return &Engine{stor: st, pl: pl, ex: exec}, nil
}

func (e *Engine) Close() error {
	return e.stor.Close()
}

// ExecSQL parses, plans and executes a single SQL statement (MVP)
func (e *Engine) ExecSQL(sql string) (any, error) {
	stmt, err := parser.Parse(sql)
	if err != nil {
		return nil, err
	}
	plan, err := e.pl.Plan(stmt)
	if err != nil {
		return nil, err
	}
	return e.ex.Execute(plan)
}
