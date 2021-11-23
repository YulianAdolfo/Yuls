package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username     = "ipsucjyl_userRegisterCH"
	password     = "RegThisPatient.21Yuls*"
	port         = "3306"
	host         = "cp-31.webhostbox.net"
	typeConn     = "tcp"
	databaseConn = "ipsucjyl_patient_clinic_history"
)

func getConnectionDB() *sql.DB {
	connection, err := sql.Open("mysql", username+":"+password+"@"+typeConn+"("+host+":"+port+")/"+databaseConn)
	if err != nil {
		fmt.Println("Error trying connection to remote database: " + err.Error())
	}
	return connection
}
