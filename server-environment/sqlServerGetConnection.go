package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
)

const (
	serverSql    = "192.168.1.198\\SQLEXPRESS"
	portSql      = 1443
	userSql      = "consentimientos"
	passwordSql  = "patientClinicH.21*"
	databaseName = "HOSVITAL"
)

func sqlServerGetConnection() *sql.DB {
	// creating the connection to sql server
	connectionQueryParams := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;", serverSql, userSql, passwordSql, portSql, databaseName)
	// creating the connection pool
	databaseConnection, err := sql.Open("sqlserver", connectionQueryParams)
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
