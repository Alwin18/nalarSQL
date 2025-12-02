package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Alwin18/nalarSQL/engine"
)

func main() {
	// Initialize engine
	e, err := engine.NewEngine(".data")
	if err != nil {
		fmt.Println("Error initializing engine:", err)
		os.Exit(1)
	}
	defer e.Close()

	// Print welcome message
	fmt.Println("┌─────────────────────────────────────┐")
	fmt.Println("│   Welcome to nalarSQL Database!    │")
	fmt.Println("│   Type 'exit' or 'quit' to exit    │")
	fmt.Println("└─────────────────────────────────────┘")
	fmt.Println()

	// Create scanner for reading input
	scanner := bufio.NewScanner(os.Stdin)

	// REPL - Read-Eval-Print Loop
	for {
		// Print prompt
		fmt.Print("nalarSQL> ")

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
			fmt.Println("Goodbye! 👋")
			break
		}

		// Execute SQL
		res, err := e.ExecSQL(input)
		if err != nil {
			fmt.Println("❌ ERROR:", err)
			continue
		}

		// Print result
		fmt.Println("✅", res)
		fmt.Println()
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
}
