#!/bin/bash

# Demo script for nalarSQL with enhanced output

# Clean up previous data
rm -rf .data

# Create SQL commands
cat << 'EOF' | ./nalarSql
CREATE TABLE employees (id INTEGER, name TEXT, department TEXT, salary INTEGER);
INSERT INTO employees (id, name, department, salary) VALUES (1, 'Alice Johnson', 'Engineering', 95000);
INSERT INTO employees (id, name, department, salary) VALUES (2, 'Bob Smith', 'Sales', 75000);
INSERT INTO employees (id, name, department, salary) VALUES (3, 'Charlie Brown', 'Engineering', 88000);
INSERT INTO employees (id, name, department, salary) VALUES (4, 'Diana Prince', 'Marketing', 82000);
SELECT * FROM employees;
UPDATE employees SET salary = 100000 WHERE id = 1;
SELECT * FROM employees;
DELETE FROM employees WHERE id = 2;
SELECT * FROM employees;
exit
EOF

