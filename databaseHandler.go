package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// DatabaseHandler is entry point to databases.
// Stores pointers to databases.
type DatabaseHandler struct {
	db *sql.DB
}

// QueryResults is alias for return type of database functions.
type QueryResults int

const (
	// RegisterOK - Register was successfull.
	RegisterOK QueryResults = 0
	// RegisterLoginExists - Login already exists.
	RegisterLoginExists QueryResults = 1
	// RegisterEmailExists - E-mail already exists.
	RegisterEmailExists QueryResults = 2
	// RegisterPasswordError - bcrypt returned an error.
	RegisterPasswordError QueryResults = 3

	// LoginOK - Login was successfull.
	LoginOK QueryResults = 0
	// LoginWrongLogin - Login is unknown or unregistered.
	LoginWrongLogin QueryResults = 1
	// LoginWrongPassword - Login is registered, but passworld doesn't match.
	LoginWrongPassword QueryResults = 2
)

// init is DatabaseHandler constructor.
/*
 - opens `SQL` database(s).
 - initializes database's tables if they aren't present

**Note:** database must be closed on program exit. */
func (p *DatabaseHandler) init(s *Server) error {
	db, err := sql.Open("sqlite3", s.srcDir+"/database.db")
	if err != nil {
		fmt.Println(err.Error())
		return StringError{"ERROR: Some of databases weren't opened!"}
	}
	p.db = db

	p.createTable()
	return nil
}

// UserRecord is single user in `userData` database.
type UserRecord struct {
	hash  string
	login string
	pass  string
	email string
	group string
	class string
}

// Creates userData table
func (p *DatabaseHandler) createTable() error {
	var table SQLiteTable
	table.setName("userData")
	table.addColumn([]SQLColumn{
		SQLiteColumn{name: "id", sqlType: "INTEGER", args: "PRIMARY KEY"},
		SQLiteColumn{name: "login", sqlType: "VARCHAR(20)", args: ""},
		SQLiteColumn{name: "password", sqlType: "VARCHAR(128)", args: ""},
		SQLiteColumn{name: "email", sqlType: "VARCHAR(128)", args: ""},
		SQLiteColumn{name: "account_type", sqlType: "VARCHAR(20)", args: ""},
		SQLiteColumn{name: "assigned_class", sqlType: "VARCHAR(20)", args: ""}})
	_, err := table.create(p.db)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

func (p *DatabaseHandler) register(login string, email string, pass string) QueryResults {
	if p.userExist(login, LoginLoginMethod) {
		return RegisterLoginExists
	}
	if p.userExist(email, EmailLoginMethod) {
		return RegisterEmailExists
	}
	passHash, err := HashPassword(pass)
	if err != nil {
		println(err.Error())
		return RegisterPasswordError
	}

	p.createRecord(UserRecord{
		hash:  "0",
		login: login,
		pass:  passHash,
		email: email,
		group: "undefined",
		class: "undefined"})
	return RegisterOK
}

func (p *DatabaseHandler) createRecord(r UserRecord) {
	register := "INSERT INTO userData " +
		"(login, password, email, account_type, assigned_class)" +
		"VALUES (?, ?, ?, ?, ?);"

	_, err := p.db.Exec(register, r.login, r.pass, r.email, r.group, r.class)
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
			"WHERE login = ?"
	} else if method == EmailLoginMethod {
		login = "SELECT email, password " +
			"FROM userData " +
			"WHERE email = ?"
	}

	result, err := p.db.Query(login, username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	for result.Next() {
		var login string
		var passHash string
		result.Scan(&login, &passHash)
		if CheckPasswordHash(pass, passHash) {
			return true
		} else {
			return false
		}
	}
	defer result.Close()
	return false
}
