package main

import (
	"database/sql"
	"fmt"
)

// SQLiteTable is implementation of `SQLTable` for `SQLite`.
type SQLiteTable struct {
	name    string
	columns []SQLColumn
}

func (p *SQLiteTable) addColumn(c []SQLColumn) {
	for _, col := range c {
		p.columns = append(p.columns, col)
	}
}

func (p *SQLiteTable) create(db *sql.DB) (string, error) {
	exec := "CREATE TABLE IF NOT EXISTS "
	exec += p.name + " ("
	for _, col := range p.columns {
		val, _ := col.get()
		exec += val + ", "
	}
	exec = exec[0:len(exec)-2] + ");"
	_, err := db.Exec(exec)
	if err != nil {
		fmt.Println(err.Error())
		return exec, err
	}
	return exec, nil
}

func (p *SQLiteTable) setName(name string) {
	p.name = name
}

// SQLiteColumn represents a `SQLite` column in `SQLTable`.
// Implements `SQLColumn`.
type SQLiteColumn struct {
	name    string
	sqlType string
	args    string
}

// get returns readable (for SQL) part of commend.
func (p SQLiteColumn) get() (string, error) {
	if p.name == "" || p.sqlType == "" {
		return "", StringError{"ERROR: SQLiteColumn is Invalid."}
	}
	return p.name + " " + p.sqlType + " " + p.args, nil
}
