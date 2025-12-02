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

// AppendRow appends a JSON-encoded row as a single line and returns the row ID
func (s *Store) AppendRow(table string, row map[string]any) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Read current rows to get next ID
	rows, err := s.scanTableUnlocked(table)
	if err != nil {
		return 0, err
	}
	nextID := int64(len(rows) + 1)

	p := s.tablePath(table)
	f, err := os.OpenFile(p, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	b, err := json.Marshal(row)
	if err != nil {
		return 0, err
	}
	if _, err := f.Write(append(b, '\n')); err != nil {
		return 0, err
	}
	return nextID, nil
}

// ScanTable naive reads and returns array of rows (as map[string]any)
func (s *Store) ScanTable(table string) ([]map[string]any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.scanTableUnlocked(table)
}

// scanTableUnlocked is the internal implementation without locking
func (s *Store) scanTableUnlocked(table string) ([]map[string]any, error) {
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

// compareValues compares two values, handling type conversions (e.g., int64 vs float64 from JSON)
func compareValues(a, b any) bool {
	if a == b {
		return true
	}

	// Handle numeric comparisons (JSON unmarshals numbers as float64)
	switch v1 := a.(type) {
	case float64:
		switch v2 := b.(type) {
		case int64:
			return v1 == float64(v2)
		case int:
			return v1 == float64(v2)
		}
	case int64:
		switch v2 := b.(type) {
		case float64:
			return float64(v1) == v2
		case int:
			return v1 == int64(v2)
		}
	case int:
		switch v2 := b.(type) {
		case float64:
			return float64(v1) == v2
		case int64:
			return int64(v1) == v2
		}
	}

	return false
}

// UpdateRows updates rows matching the where condition and returns count of updated rows
func (s *Store) UpdateRows(table string, set map[string]any, whereCol string, whereVal any) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.scanTableUnlocked(table)
	if err != nil {
		return 0, err
	}

	updated := 0
	for _, row := range rows {
		// If no where clause, update all rows
		if whereCol == "" {
			for k, v := range set {
				row[k] = v
			}
			updated++
		} else {
			// Check if row matches where condition
			if rowVal, ok := row[whereCol]; ok && compareValues(rowVal, whereVal) {
				for k, v := range set {
					row[k] = v
				}
				updated++
			}
		}
	}

	// Rewrite the entire table file
	if err := s.rewriteTable(table, rows); err != nil {
		return 0, err
	}

	return updated, nil
}

// DeleteRows deletes rows matching the where condition and returns count of deleted rows
func (s *Store) DeleteRows(table string, whereCol string, whereVal any) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.scanTableUnlocked(table)
	if err != nil {
		return 0, err
	}

	var newRows []map[string]any
	deleted := 0

	for _, row := range rows {
		// If no where clause, delete all rows
		if whereCol == "" {
			deleted++
			continue
		}
		// Check if row matches where condition
		if rowVal, ok := row[whereCol]; ok && compareValues(rowVal, whereVal) {
			deleted++
			continue
		}
		newRows = append(newRows, row)
	}

	// Rewrite the entire table file
	if err := s.rewriteTable(table, newRows); err != nil {
		return 0, err
	}

	return deleted, nil
}

// rewriteTable rewrites the entire table file with new rows
func (s *Store) rewriteTable(table string, rows []map[string]any) error {
	p := s.tablePath(table)

	// Read the header first
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(f)
	var header map[string]any
	if err := dec.Decode(&header); err != nil {
		f.Close()
		return err
	}
	f.Close()

	// Write to temp file
	tmpPath := p + ".tmp"
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(tmpFile)
	// Write header
	if err := enc.Encode(header); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return err
	}

	// Write all rows
	for _, row := range rows {
		b, err := json.Marshal(row)
		if err != nil {
			tmpFile.Close()
			os.Remove(tmpPath)
			return err
		}
		if _, err := tmpFile.Write(append(b, '\n')); err != nil {
			tmpFile.Close()
			os.Remove(tmpPath)
			return err
		}
	}

	tmpFile.Close()

	// Replace original file with temp file
	return os.Rename(tmpPath, p)
}
