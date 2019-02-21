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
	RegisterOK  RegisterErrors = 0
	LoginExists RegisterErrors = 1
	EmailExists RegisterErrors = 2
)

func (p *DatabaseHandler) init() {
	db, err := sql.Open("sqlite3", "./src/database.db")
	p.db = db
	if err != nil {
		fmt.Println(err.Error())
	}
	p.createTable() // should be separate program
}

func (p *DatabaseHandler) createTable() bool {
	createTable :=
		"CREATE TABLE IF NOT EXISTS userData(" +
			"id INT(10)," +
			"login VARCHAR(20)," +
			"password VARCHAR(128)," +
			"email VARCHAR(128)," +
			"account_type VARCHAR(20)," +
			"assigned_class VARCHAR(20));"

	_, err := p.db.Exec(createTable)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func (p *DatabaseHandler) register(login string, email string, pass string) RegisterErrors {
	if p.userExist(login, LoginLoginMethod) {
		return LoginExists
	}
	if p.userExist(email, EmailLoginMethod) {
		return EmailExists
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
