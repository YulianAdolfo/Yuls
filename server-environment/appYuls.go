package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"
)

// fields in database
type dataPatientHC struct {
	ActualDateRegistry string
	DateClinicHistory  string
	IdPatient          int
	PatientNames       string
	PatientLastnames   string
	TypeId             string
	HasError           bool
}
type returnMessage struct {
	ContenMessage string
}
type sqlColumnsName struct {
	MPNom1, MPNom2, MPApe1, MPApe2, MPTDoc string
}
type resultPatientSqlServer struct {
	FirstName, SecondName, FirstLastname, SecondLastname, TypId string
}

const (
	DATABASE_IN_USE = "clinic_history_Test"
)

// insert new patients
func newClinicHistory(dataPatienStruct dataPatientHC) error {
	// connecting to database
	connection := getConnectionDB()
	insertQuery := "INSERT INTO " + DATABASE_IN_USE + " (actualDateRegistry, dateClinicHistory, IdPatient, patientNames, patientLastnames, typeId, hasError) VALUES (?,?,?,?,?,?,?)"
	contextQuery, cancelFunction := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunction()
	statement, err := connection.PrepareContext(contextQuery, insertQuery)
	if err != nil {
		fmt.Println("Error preparing the query with context " + err.Error())
		return err
	}
	defer statement.Close()
	_, err = statement.ExecContext(contextQuery, dataPatienStruct.ActualDateRegistry, dataPatienStruct.DateClinicHistory, dataPatienStruct.IdPatient, dataPatienStruct.PatientNames, dataPatienStruct.PatientLastnames, dataPatienStruct.TypeId, dataPatienStruct.HasError)
	if err != nil {
		fmt.Println("Error executing the context " + err.Error())
		return err
	}
	defer connection.Close()
	fmt.Println("success")

	return err
}
func setPatientRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		bodyRequest, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error reading the body: " + err.Error())
		}
		var dataPatientHc dataPatientHC
		json.Unmarshal(bodyRequest, &dataPatientHc)
		// getting the actual date
		dataPatientHc.ActualDateRegistry = time.Now().String()
		// function to insert new records
		err = newClinicHistory(dataPatientHc)
		if err != nil {
			fmt.Fprint(w, responseClientError(err))
			return
		}
		fmt.Println(responseClientSucess())
		fmt.Fprint(w, responseClientSucess())
	}
}
func responseClientError(err error) string {
	m := returnMessage{
		ContenMessage: err.Error(),
	}
	contentJson, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(contentJson)
}
func responseClientSucess() string {
	m := returnMessage{
		ContenMessage: "successfull",
	}
	contentJson, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(contentJson)
}
func getInfoPatientFromHosvital(id string, connectionSqlServer *sql.DB) string {
	contextConnection := context.Background()
	// check if the connection is alive
	err := connectionSqlServer.PingContext(contextConnection)
	if err != nil {
		fmt.Println("Error in ping connection to sqlserver: " + err.Error())
	}
	// if the connection is alive so create the sql qery

	sqlGetInfo := fmt.Sprintf("SELECT MPNom1, MPNom2, MPApe1, MPApe2, MPTDoc FROM CAPBAS WHERE MPCedu =" + "'" + id + "'")
	fmt.Println(sqlGetInfo)
	rows, err := connectionSqlServer.QueryContext(contextConnection, sqlGetInfo)
	if err != nil {
		fmt.Println("Error executing the context to to sql server: " + err.Error())
	}
	defer rows.Close()

	var dataResultPatientHosvital resultPatientSqlServer
	for rows.Next() {
		var dataSqlServer sqlColumnsName
		err = rows.Scan(&dataSqlServer.MPNom1, &dataSqlServer.MPNom2, &dataSqlServer.MPApe1, &dataSqlServer.MPApe2, &dataSqlServer.MPTDoc)
		if err != nil {
			fmt.Println("Error scannig data from sql-server: " + err.Error())
		}
		// asigning values to then convert then into json
		dataResultPatientHosvital.FirstName = dataSqlServer.MPNom1
		dataResultPatientHosvital.SecondName = dataSqlServer.MPNom2
		dataResultPatientHosvital.FirstLastname = dataSqlServer.MPApe1
		dataResultPatientHosvital.SecondLastname = dataSqlServer.MPApe2
		dataResultPatientHosvital.TypId = dataSqlServer.MPTDoc
	}
	defer connectionSqlServer.Close()
	toJsonData, err := json.Marshal(dataResultPatientHosvital)
	if err != nil {
		fmt.Println("Error marshalling the json: " + err.Error())
	}
	return string(toJsonData)
}
func patientHosvital(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		idPatient := r.URL.Query().Get("id-patient")
		data := getInfoPatientFromHosvital(idPatient, sqlServerGetConnection())
		fmt.Fprint(w, data)
	}
}
func app(w http.ResponseWriter, r *http.Request) {
	appTemplate := template.Must(template.ParseFiles("../client-environment/app.html"))
	appTemplate.Execute(w, nil)
}
func main() {
	publicElementsApp := http.FileServer(http.Dir("../public"))
	http.Handle("/public/", http.StripPrefix("/public/", publicElementsApp))
	fmt.Println("Using the database: " + DATABASE_IN_USE)
	http.HandleFunc("/record-patient", setPatientRecord)
	http.HandleFunc("/get-data-patient", patientHosvital)
	http.HandleFunc("/Yuls", app)
	http.ListenAndServe(":8005", nil)
}
