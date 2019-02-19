package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)




func connectDatabase(){
	db, err := sql.Open("mysql", "username:password@tcp(host:3306)/ejournal")

	if checkError(err){
		fmt.Println("Wystapil blad poczas laczenia z baza danych.")
	}else {
		fmt.Println("Pomyslnie polaczono z bazą danych")
	}


	statement, err := db.Query("CREATE TABLE IF NOT EXISTS `userData`(`id` INT(10) NOT NULL AUTO_INCREMENT, `username` VARCHAR(20), password VARCHAR(128), email VARCHAR(128), account_type VARCHAR(20), assigned_class VARCHAR(20), PRIMARY KEY(id))")
	if checkError(err) == false{
		fmt.Println("Pomyslnie załadowano tabelę userData")
	}
	defer statement.Close()
	defer db.Close()
}


func checkError(err error) bool{
	if err != nil {
		panic(err.Error())
		return true
	}
	return false
}