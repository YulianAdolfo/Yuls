package main

import (
	"database/sql"
	"fmt"

	"Yuls/readerparams"

	_ "github.com/go-sql-driver/mysql"
)

func getConnectionDB() *sql.DB {
	username, password, typeConn, host, port, databaseConn := readerparams.ReadConnectionMySqlParameters()
	connection, err := sql.Open("mysql", username+":"+password+"@"+typeConn+"("+host+":"+port+")/"+databaseConn)
	if err != nil {
		fmt.Println("Error trying connection to remote database: " + err.Error())
	}
	err = connection.Ping()
	if err != nil {
		fmt.Println("Error in ping connection: " + err.Error())
	}
	return connection
}
