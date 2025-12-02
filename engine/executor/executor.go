package executor

import (
	"fmt"

	"github.com/Alwin18/nalarSQL/engine/planner"
	"github.com/Alwin18/nalarSQL/engine/storage"
)

type Executor struct {
	store *storage.Store
}

func NewExecutor(store *storage.Store) *Executor {
	return &Executor{store: store}
}

func (e *Executor) Execute(plan planner.Plan) (any, error) {
	switch p := plan.(type) {
	case *planner.PlanCreateTable:
		cols := make([]storage.ColumnDefinition, len(p.Stmt.Columns))
		for i, c := range p.Stmt.Columns {
			cols[i] = storage.ColumnDefinition{Name: c.Name, Type: c.Type}
		}
		return nil, e.store.CreateTable(p.Stmt.TableName, cols)
	case *planner.PlanInsert:
		// TODO: map values to columns
		row := map[string]any{"id": p.Stmt.Values[0], "name": p.Stmt.Values[1], "age": p.Stmt.Values[2]}
		if err := e.store.AppendRow(p.Stmt.Table, row); err != nil {
			return nil, err
		}
		return "OK", nil
	case *planner.PlanSelect:
		rows, err := e.store.ScanTable(p.Stmt.Table)
		if err != nil {
			return nil, err
		}
		return rows, nil
	default:
		return nil, fmt.Errorf("executor: unsupported plan type %T", p)
	}
}
