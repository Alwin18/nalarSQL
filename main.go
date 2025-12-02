package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Alwin18/nalarSQL/engine"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

func main() {
	// Initialize engine
	e, err := engine.NewEngine(".data")
	if err != nil {
		fmt.Printf("%s❌ Error initializing engine: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}
	defer e.Close()

	// Print welcome message
	fmt.Println(colorCyan + "┌─────────────────────────────────────┐" + colorReset)
	fmt.Println(colorCyan + "│" + colorBold + "   Welcome to nalarSQL Database!    " + colorReset + colorCyan + "│" + colorReset)
	fmt.Println(colorCyan + "│" + colorReset + "   Type 'exit' or 'quit' to exit    " + colorCyan + "│" + colorReset)
	fmt.Println(colorCyan + "└─────────────────────────────────────┘" + colorReset)
	fmt.Println()

	// Create scanner for reading input
	scanner := bufio.NewScanner(os.Stdin)

	// REPL - Read-Eval-Print Loop
	for {
		// Print prompt
		fmt.Print(colorGreen + "nalarSQL> " + colorReset)

		// Read input
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if input == "" {
			continue
		}

		// Check for exit commands
		inputLower := strings.ToLower(input)
		if inputLower == "exit" || inputLower == "quit" || inputLower == "exit;" || inputLower == "quit;" {
			fmt.Println(colorYellow + "👋 Goodbye!" + colorReset)
			break
		}

		// Execute SQL
		res, err := e.ExecSQL(input)
		if err != nil {
			fmt.Printf("%s❌ ERROR: %v%s\n\n", colorRed, err, colorReset)
			continue
		}

		// Print formatted result
		printResult(res)
		fmt.Println()
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Printf("%s❌ Error reading input: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}
}

// printResult formats and prints the query result in a user-friendly way
func printResult(result any) {
	if result == nil {
		fmt.Printf("%s✅ Query executed successfully%s\n", colorGreen, colorReset)
		return
	}

	switch v := result.(type) {
	case []map[string]any:
		// SELECT result - print as table
		printTable(v)
	case map[string]any:
		// INSERT/UPDATE/DELETE result
		printOperationResult(v)
	default:
		fmt.Printf("%s✅ %v%s\n", colorGreen, result, colorReset)
	}
}

// printTable prints a slice of maps as a formatted table
func printTable(rows []map[string]any) {
	if len(rows) == 0 {
		fmt.Printf("%s📭 No rows returned%s\n", colorYellow, colorReset)
		return
	}

	// Get all column names
	colMap := make(map[string]bool)
	for _, row := range rows {
		for col := range row {
			colMap[col] = true
		}
	}

	// Convert to sorted slice for consistent ordering
	var columns []string
	for col := range colMap {
		columns = append(columns, col)
	}

	// Calculate column widths
	widths := make(map[string]int)
	for _, col := range columns {
		widths[col] = len(col)
	}
	for _, row := range rows {
		for _, col := range columns {
			valStr := fmt.Sprintf("%v", row[col])
			if len(valStr) > widths[col] {
				widths[col] = len(valStr)
			}
		}
	}

	// Print top border
	fmt.Print(colorBlue + "┌")
	for i, col := range columns {
		fmt.Print(strings.Repeat("─", widths[col]+2))
		if i < len(columns)-1 {
			fmt.Print("┬")
		}
	}
	fmt.Println("┐" + colorReset)

	// Print header
	fmt.Print(colorBlue + "│" + colorReset)
	for _, col := range columns {
		fmt.Printf(" %s%-*s%s ", colorBold+colorCyan, widths[col], col, colorReset)
		fmt.Print(colorBlue + "│" + colorReset)
	}
	fmt.Println()

	// Print separator
	fmt.Print(colorBlue + "├")
	for i, col := range columns {
		fmt.Print(strings.Repeat("─", widths[col]+2))
		if i < len(columns)-1 {
			fmt.Print("┼")
		}
	}
	fmt.Println("┤" + colorReset)

	// Print rows
	for _, row := range rows {
		fmt.Print(colorBlue + "│" + colorReset)
		for _, col := range columns {
			val := row[col]
			valStr := fmt.Sprintf("%v", val)
			fmt.Printf(" %-*s ", widths[col], valStr)
			fmt.Print(colorBlue + "│" + colorReset)
		}
		fmt.Println()
	}

	// Print bottom border
	fmt.Print(colorBlue + "└")
	for i, col := range columns {
		fmt.Print(strings.Repeat("─", widths[col]+2))
		if i < len(columns)-1 {
			fmt.Print("┴")
		}
	}
	fmt.Println("┘" + colorReset)

	// Print row count
	rowCount := len(rows)
	rowWord := "row"
	if rowCount != 1 {
		rowWord = "rows"
	}
	fmt.Printf("%s%d %s returned%s\n", colorGray, rowCount, rowWord, colorReset)
}

// printOperationResult prints results from INSERT/UPDATE/DELETE operations
func printOperationResult(result map[string]any) {
	if rowid, ok := result["rowid"]; ok {
		fmt.Printf("%s✅ Row inserted with ID: %v%s\n", colorGreen, rowid, colorReset)
		return
	}

	if updated, ok := result["updated"]; ok {
		count := updated.(int)
		rowWord := "row"
		if count != 1 {
			rowWord = "rows"
		}
		if count > 0 {
			fmt.Printf("%s✅ %d %s updated%s\n", colorGreen, count, rowWord, colorReset)
		} else {
			fmt.Printf("%s⚠️  No rows matched the WHERE condition%s\n", colorYellow, colorReset)
		}
		return
	}

	if deleted, ok := result["deleted"]; ok {
		count := deleted.(int)
		rowWord := "row"
		if count != 1 {
			rowWord = "rows"
		}
		if count > 0 {
			fmt.Printf("%s✅ %d %s deleted%s\n", colorGreen, count, rowWord, colorReset)
		} else {
			fmt.Printf("%s⚠️  No rows matched the WHERE condition%s\n", colorYellow, colorReset)
		}
		return
	}

	// Fallback for unknown result format
	fmt.Printf("%s✅ %v%s\n", colorGreen, result, colorReset)
}
