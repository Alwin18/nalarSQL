# Nalar SQL

Nalar SQL is a simple SQL database management system written in Go.

```
/engine
    parser/
        lexer.go
        parser.go
        ast.go
    planner/
        planner.go
    executor/
        insert.go
        update.go
        delete.go
        select.go
    storage/
        file.go
        page.go
        row.go
        index.go
    utils/
        errors.go

main.go
```