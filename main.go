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
		`INSERT INTO users (id, name, age) VALUES (3, 'Charlie', 35);`,
		`SELECT * FROM users;`,
		`UPDATE users SET age = 31 WHERE id = 1;`,
		`SELECT * FROM users;`,
		`DELETE FROM users WHERE id = 2;`,
		`SELECT * FROM users;`,
	}

	for i, s := range sqls {
		fmt.Printf("\n[%d] %s\n", i+1, s)
		res, err := e.ExecSQL(s)
		if err != nil {
			fmt.Println("ERROR:", err)
			continue
		}
		fmt.Println("=>", res)
	}
}
