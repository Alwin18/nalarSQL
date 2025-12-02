package main

import (
	"fmt"

	"github.com/Alwin18/nalarSQL/engine"
)

func main() {
	// Quick startup: initialize engine and run a sample SQL
	e, err := engine.NewEngine(".data")
	if err != nil {
		panic(err)
	}
	defer e.Close()

	sqls := []string{
		`CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, age INTEGER);`,
		`INSERT INTO users (id, name, age) VALUES (1, 'Alice', 30);`,
		`INSERT INTO users (id, name, age) VALUES (2, 'Bob', 25);`,
		`SELECT id, name, age FROM users;`,
	}

	for _, s := range sqls {
		res, err := e.ExecSQL(s)
		if err != nil {
			fmt.Println("ERROR:", err)
			continue
		}
		fmt.Println("=>", res)
	}
}
