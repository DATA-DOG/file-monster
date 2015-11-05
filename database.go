package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func openSQLConnection() {
	conn, err := sql.Open("mysql", fmt.Sprint(config.DbUsername, ":", config.DbPassword, "@/", config.DbDatabase))
	checkErr(err)

	db = conn

	log.Print("Connected to SQL")
}

func dbGetUserToken(user string) string {
	stmt, err := db.Prepare("SELECT token FROM access_token WHERE user = ?")
	checkErr(err)

	row := stmt.QueryRow(user)

	var token string
	row.Scan(&token)

	return token
}

func dbSaveUserToken(user string, token string) {
	stmt, err := db.Prepare("INSERT INTO access_token (user, token) VALUES (?, ?)")
	checkErr(err)

	_, err = stmt.Exec(user, token)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
