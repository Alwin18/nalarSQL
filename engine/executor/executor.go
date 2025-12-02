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
		row := map[string]any{}
		for i, c := range p.Stmt.Columns {
			if i < len(p.Stmt.Values) {
				row[c] = p.Stmt.Values[i]
			}
		}
		id, err := e.store.AppendRow(p.Stmt.Table, row)
		if err != nil {
			return nil, err
		}
		return map[string]any{"rowid": id}, nil
	case *planner.PlanSelect:
		rows, err := e.store.ScanTable(p.Stmt.Table)
		if err != nil {
			return nil, err
		}
		// naive: ignore projection, return all columns
		return rows, nil
	case *planner.PlanUpdate:
		// apply update: simple where only supports equality on a column (col = value)
		updated, err := e.store.UpdateRows(p.Stmt.Table, p.Stmt.Set, p.Stmt.WhereColumn, p.Stmt.WhereValue)
		if err != nil {
			return nil, err
		}
		return map[string]any{"updated": updated}, nil
	case *planner.PlanDelete:
		deleted, err := e.store.DeleteRows(p.Stmt.Table, p.Stmt.WhereColumn, p.Stmt.WhereValue)
		if err != nil {
			return nil, err
		}
		return map[string]any{"deleted": deleted}, nil
	default:
		return nil, fmt.Errorf("executor: unsupported plan type %T", p)
	}
}
