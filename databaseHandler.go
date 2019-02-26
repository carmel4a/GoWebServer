package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DatabaseHandler struct {
	db *sql.DB
}

type RegisterErrors int

const (
	RegisterOK          RegisterErrors = 0
	RegisterLoginExists RegisterErrors = 1
	RegisterEmailExists RegisterErrors = 2
)

func (p *DatabaseHandler) init(s *Server) {
	db, err := sql.Open("sqlite3", s.srcDir+"/database.db")
	p.db = db
	if err != nil {
		fmt.Println(err.Error())
	}
	p.createTable()
}

type Column struct {
	name    string
	sqlType string
	args    string
}

func (p Column) get() string {
	return p.name + " " + p.sqlType + " " + p.args
}

type SQLTable struct {
	name    string
	columns []Column
}

func (p *SQLTable) addColumn(c []Column) {
	for _, col := range c {
		p.columns = append(p.columns, col)
	}
}

func (p *SQLTable) create(db *sql.DB) (string, error) {
	exec := "CREATE TABLE IF NOT EXISTS "
	exec += p.name + " ("
	for _, col := range p.columns {
		exec += col.get() + ", "
	}
	exec = exec[0:len(exec)-2] + ");"
	_, err := db.Exec(exec)
	if err != nil {
		fmt.Println(err.Error())
		return exec, err
	}
	return exec, nil
}

func (p *DatabaseHandler) createTable() error {
	table := SQLTable{name: "userData"}
	table.addColumn([]Column{
		Column{name: "id", sqlType: "INT(10)", args: ""},
		Column{name: "login", sqlType: "VARCHAR(20)", args: ""},
		Column{name: "password", sqlType: "VARCHAR(128)", args: ""},
		Column{name: "email", sqlType: "VARCHAR(128)", args: ""},
		Column{name: "account_type", sqlType: "VARCHAR(20)", args: ""},
		Column{name: "assigned_class", sqlType: "VARCHAR(20)", args: ""}})
	_, err := table.create(p.db)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

func (p *DatabaseHandler) register(login string, email string, pass string) RegisterErrors {
	if p.userExist(login, LoginLoginMethod) {
		return RegisterLoginExists
	}
	if p.userExist(email, EmailLoginMethod) {
		return RegisterEmailExists
	}
	p.createRecord("0", login, pass, email, "undefined", "undefined")
	return RegisterOK
}

func (p *DatabaseHandler) createRecord(hash string, login string, pass string, email string, group string, class string) {

	register := "INSERT INTO userData " +
		"(id, login, password, email, account_type, assigned_class)" +
		"VALUES (?, ?, ?, ?, ?, ?);"

	_, err := p.db.Exec(register, hash, login, pass, email, group, class)
	if err != nil {
		fmt.Println(err.Error())
	}
}

type LoginMethod string

const (
	LoginLoginMethod LoginMethod = "login"
	EmailLoginMethod LoginMethod = "email"
)

func (p *DatabaseHandler) userExist(username string, method LoginMethod) bool {
	userExist := "SELECT ? " +
		"FROM userData " +
		"WHERE ? = ?;"

	result, err := p.db.Query(userExist, method, method, username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	for result.Next() {
		return true
	}
	defer result.Close()
	return false
}

func (p *DatabaseHandler) login(username string, pass string, method LoginMethod) bool {
	var login string
	if method == LoginLoginMethod {
		login = "SELECT login, password " +
			"FROM userData " +
			"WHERE login = ? AND password = ?"
	} else if method == EmailLoginMethod {
		login = "SELECT email, password " +
			"FROM userData " +
			"WHERE email = ? AND password = ?"
	}

	result, err := p.db.Query(login, username, pass)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	for result.Next() {
		return true
	}
	defer result.Close()
	return false
}
