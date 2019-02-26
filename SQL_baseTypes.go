package main

import "database/sql"

// SQLTable is base for `SQL` tables.
type SQLTable interface {
	// Append column.
	addColumn([]SQLColumn)
	// Creates full table in database, Do not override existing data.
	create(db *sql.DB) (string, error)
	setName(string)
}

// SQLColumn is base interface for `SQL` columns.
type SQLColumn interface {
	// get returns readable (for SQL) part of commend to create a column.
	get() (string, error)
}
