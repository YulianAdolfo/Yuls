package main

import (
	"context"
	"encoding/json"
	"fmt"
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
func newClinicHistory(data []byte) {
	var dataPatienStruct dataPatientHC
	json.Unmarshal(data, &dataPatienStruct)
	fmt.Println(dataPatienStruct.ActualDateRegistry)
	// connecting to database
	connection := getConnectionDB()
	insertQuery := "INSERT INTO " + DATABASE_IN_USE + " (actualDateRegistry, dateClinicHistory, IdPatient, patientNames, patientLastnames, typeId, hasError) VALUES (?,?,?,?,?,?,?)"
	contextQuery, cancelFunction := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunction()
	statement, err := connection.PrepareContext(contextQuery, insertQuery)
	if err != nil {
		fmt.Println("Error preparing the query with context " + err.Error())
		return
	}
	defer statement.Close()
	_, err = statement.ExecContext(contextQuery, dataPatienStruct.ActualDateRegistry, dataPatienStruct.DateClinicHistory, dataPatienStruct.IdPatient, dataPatienStruct.PatientNames, dataPatienStruct.PatientLastnames, dataPatienStruct.TypeId, dataPatienStruct.HasError)
	if err != nil {
		fmt.Println("Error executing the context " + err.Error())
		return
	}
	defer connection.Close()

}
func registry() {
	data := getContentBody()
	newClinicHistory(data)
}
func getContentBody() []byte {
	var dataPatientHc dataPatientHC
	dataPatientHc.ActualDateRegistry = time.Now().String()
	dataPatientHc.DateClinicHistory = "2021-10-12"
	dataPatientHc.PatientNames = "Yulian Adolfo"
	dataPatientHc.PatientLastnames = "Rojas Gañan"
	dataPatientHc.TypeId = "CC"
	dataPatientHc.IdPatient = 1007441849
	dataPatientHc.HasError = true

	jsonContent, err := json.Marshal(dataPatientHc)
	if err != nil {
		fmt.Println("Error marshalling in getContentBody: " + err.Error())
	}
	return jsonContent
}
func main() {
	fmt.Println("Using the database: " + DATABASE_IN_USE)
	registry()
}
