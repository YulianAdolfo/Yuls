package main

import (
	"context"
	"database/sql"
	"fmt"

	"Yuls/readerparams"

	_ "github.com/denisenkom/go-mssqldb"
)

func sqlServerGetConnection() *sql.DB {
	// getting parameters from the json file
	serverSql, userSql, passwordSql, portSql, databaseName := readerparams.ReadConnectionSqlParameters()
	// creating the connection to sql server
	connectionQueryParams := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", serverSql, userSql, passwordSql, portSql, databaseName)
	// creating the connection pool
	databaseConnection, err := sql.Open("mssql", connectionQueryParams)
	if err != nil {
		fmt.Println("Error opening the connection to sql server: " + err.Error())
	}
	ctx := context.Background()
	// Ping database to see if it's still alive.
	// Important for handling network issues and long queries.
	err = databaseConnection.PingContext(ctx)
	if err != nil {
		fmt.Println("Error pinging database: " + err.Error())
	}
	return databaseConnection
}
