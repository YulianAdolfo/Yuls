package main

import (
	"context"
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
		dataPatientHc.ActualDateRegistry = time.Now().String() // getting the actual date
		err = newClinicHistory(dataPatientHc)                  // function to insert new records

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
	http.HandleFunc("/Yuls", app)
	http.ListenAndServe(":8005", nil)
}
