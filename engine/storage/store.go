package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type Store struct {
	baseDir string
	mu      sync.RWMutex
}

func NewStore(baseDir string) (*Store, error) {
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, err
	}
	return &Store{baseDir: filepath.Clean(baseDir)}, nil
}

func (s *Store) Close() error {
	return nil
}

func (s *Store) tablePath(name string) string {
	return filepath.Join(s.baseDir, name+".tbl")
}

// ColumnDefinition represents a column in the table schema
type ColumnDefinition struct {
	Name string
	Type string
}

// CreateTable writes a schema file (very simple JSON header)
func (s *Store) CreateTable(name string, cols []ColumnDefinition) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p := s.tablePath(name)
	f, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("table %s already exists", name)
		}
		return err
	}
	defer f.Close()

	head := map[string]any{"columns": cols}

	enc := json.NewEncoder(f)
	if err := enc.Encode(head); err != nil {
		return err
	}

	return nil
}

// AppendRow appends a JSON-encoded row as a single line
func (s *Store) AppendRow(table string, row map[string]any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p := s.tablePath(table)
	f, err := os.OpenFile(p, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.Marshal(row)
	if err != nil {
		return err
	}
	if _, err := f.Write(append(b, '\n')); err != nil {
		return err
	}
	return nil
}

// ScanTable naive reads and returns array of rows (as map[string]any)
func (s *Store) ScanTable(table string) ([]map[string]any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p := s.tablePath(table)
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	// First token is header
	var header map[string]any
	if err := dec.Decode(&header); err != nil {
		return nil, err
	}
	rows := make([]map[string]any, 0)
	for {
		var m map[string]any
		if err := dec.Decode(&m); err != nil {
			// EOF expected
			if errors.Is(err, io.EOF) || err.Error() == "EOF" {
				break
			}
			break
		}
		rows = append(rows, m)
	}
	return rows, nil
}
