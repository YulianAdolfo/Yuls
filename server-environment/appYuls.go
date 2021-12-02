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
	HasError           string
}
type returnMessage struct {
	ContenMessage string
}
type sqlColumnsName struct {
	MPNom1, MPApe1, MPTDoc string
}
type resultPatientSqlServer struct {
	Names, Lastnames, TypId string
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
func patientHosvital(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		idPatient := r.URL.Query().Get("id-patient")
		data := getInfoPatientFromHosvitalTest(idPatient, sqlServerGetConnection())
		fmt.Fprint(w, data)
	}
}
func app(w http.ResponseWriter, r *http.Request) {
	appTemplate := template.Must(template.ParseFiles("../client-environment/app.html"))
	appTemplate.Execute(w, nil)
}
func getInfoPatientFromHosvitalTest(id string, connectionSqlServer *sql.DB) string {
	fmt.Println(id)
	contextConnection := context.Background()
	// check if the connection is alive
	err := connectionSqlServer.PingContext(contextConnection)
	if err != nil {
		fmt.Println("Error in ping connection to sqlserver: " + err.Error())
	}
	// if the connection is alive so create the sql qery

	sqlGetInfo := fmt.Sprintf("SELECT\n" +
		"RTRIM(CONCAT(CONCAT(LEFT(MPNom1, 1), LOWER(RIGHT(RTRIM(MPNom1), LEN(MPNom1)-1))),' ',IIF (LEN(RTRIM(MPNom2))=0,'', CONCAT(LEFT(MPNom2, 1), LOWER(RIGHT(RTRIM(MPNom2), LEN(MPNom2)-1)))))),\n+" +
		"RTRIM(CONCAT(CONCAT(LEFT(MPApe1, 1), LOWER(RIGHT(RTRIM(MPApe1), LEN(MPApe1)-1))),' ',IIF (LEN(RTRIM(MPApe2))=0,'', CONCAT(LEFT(MPApe2, 1), LOWER(RIGHT(RTRIM(MPApe2), LEN(MPApe2)-1)))))),\n" +
		"\n" +
		"CASE\n" +
		"WHEN MPTDoc = 'CC'  THEN 0 	\n" +
		"WHEN MPTDoc = 'TI'  THEN 1 	\n" +
		"WHEN MPTDoc = 'CE'  THEN 2 	\n" +
		"WHEN MPTDoc = 'ASI' THEN 3 	\n" +
		"WHEN MPTDoc = 'CI'  THEN 4 	\n" +
		"WHEN MPTDoc = 'MSI' THEN 5 	\n" +
		"WHEN MPTDoc = 'NU'  THEN 6 	\n" +
		"WHEN MPTDoc = 'PA'  THEN 7 	\n" +
		"WHEN MPTDoc = 'PE'  THEN 8 	\n" +
		"WHEN MPTDoc = 'RC'  THEN 9 	\n" +
		"WHEN MPTDoc = 'RI'  THEN 10 	\n" +
		"WHEN MPTDoc = 'PEP' THEN 11 	\n" +
		"WHEN MPTDoc = 'NIT' THEN 12	\n" +
		"ELSE 0 						\n" +
		"END							\n" +
		"\n" +
		"FROM CAPBAS WHERE MPCedu =" + "'" + id + "'")
	rows, err := connectionSqlServer.QueryContext(contextConnection, sqlGetInfo)
	if err != nil {
		fmt.Println("Error executing the context to sql server: " + err.Error())
	}
	defer rows.Close()

	var dataResultPatientHosvital resultPatientSqlServer
	for rows.Next() {
		var dataSqlServer sqlColumnsName
		err = rows.Scan(&dataSqlServer.MPNom1, &dataSqlServer.MPApe1, &dataSqlServer.MPTDoc)
		if err != nil {
			fmt.Println("Error scannig data from sql-server: " + err.Error())
		}
		// asigning values to then convert then into json
		dataResultPatientHosvital.Names = dataSqlServer.MPNom1
		dataResultPatientHosvital.Lastnames = dataSqlServer.MPApe1
		dataResultPatientHosvital.TypId = dataSqlServer.MPTDoc
	}
	defer connectionSqlServer.Close()
	toJsonData, err := json.Marshal(dataResultPatientHosvital)
	if err != nil {
		fmt.Println("Error marshalling the json: " + err.Error())
	}
	fmt.Println(string(toJsonData))
	return string(toJsonData)
}
func selectingDataToBuildReport(dateStart, dateEnd, typeDateReport, showOnlyErrors string) ([]string, error) {
	connection := getConnectionDB()
	query, err := connection.Query("CALL GET_REPORT(?, ?, ?, ?)", typeDateReport, dateStart, dateEnd, showOnlyErrors)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	var informationForReport []string
	for query.Next() {
		var dataReport dataPatientHC
		err = query.Scan(&dataReport.IdPatient, &dataReport.TypeId, &dataReport.DateClinicHistory, &dataReport.ActualDateRegistry, &dataReport.PatientNames, &dataReport.PatientLastnames, &dataReport.HasError)
		if err != nil {
			fmt.Println("Error scannig : " + err.Error())
		}
		content, err := json.Marshal(dataPatientHC{
			IdPatient:          dataReport.IdPatient,
			TypeId:             dataReport.TypeId,
			DateClinicHistory:  dataReport.DateClinicHistory,
			ActualDateRegistry: dataReport.ActualDateRegistry,
			PatientNames:       dataReport.PatientNames,
			PatientLastnames:   dataReport.PatientLastnames,
			HasError:           dataReport.HasError,
		})
		if err != nil {
			fmt.Println("Error marshalling data: " + err.Error())
		}
		informationForReport = append(informationForReport, string(content))
	}
	defer connection.Close()
	return informationForReport, nil
}
func getReport(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		dateStart := r.URL.Query().Get("date-start")
		dateEnd := r.URL.Query().Get("date-end")
		checkPatientErrors := r.URL.Query().Get("check-only-p-errors")
		generateReportBy := r.URL.Query().Get("gen-by")

		dataRequest, err := selectingDataToBuildReport(dateStart, dateEnd, generateReportBy, checkPatientErrors)
		if err != nil {
			fmt.Println("Error: " + err.Error())
		}
		toJson, err := json.Marshal(dataRequest)
		if err != nil {
			fmt.Println("Error marshalling: " + err.Error())
		}
		fmt.Fprint(w, string(toJson))
	}
}
func main() {
	publicElementsApp := http.FileServer(http.Dir("../public"))
	http.Handle("/public/", http.StripPrefix("/public/", publicElementsApp))
	fmt.Println("Using the database: " + DATABASE_IN_USE)
	http.HandleFunc("/record-patient", setPatientRecord)
	http.HandleFunc("/get-data-patient", patientHosvital)
	http.HandleFunc("/get-information-from-patient", getReport)
	http.HandleFunc("/Yuls", app)
	http.ListenAndServe(":8005", nil)
}
