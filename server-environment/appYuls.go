package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
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
type setDataExcel struct {
	DataExcel []string
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
		ContenMessage: "successful",
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
		data, err := getInfoPatientFromHosvitalTest(idPatient, sqlServerGetConnection())
		if err != nil {
			fmt.Fprint(w, responseClientError(err))
		}
		fmt.Fprint(w, data)
	}
}
func app(w http.ResponseWriter, r *http.Request) {
	appTemplate := template.Must(template.ParseFiles("../client-environment/app.html"))
	appTemplate.Execute(w, nil)
}
func getInfoPatientFromHosvitalTest(id string, connectionSqlServer *sql.DB) (string, error) {
	contextConnection := context.Background()
	// check if the connection is alive
	err := connectionSqlServer.PingContext(contextConnection)
	if err != nil {
		fmt.Println("Error in ping connection to Hosvital " + err.Error())
		return "", err

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
		fmt.Println("Error connection to Hosvital: " + err.Error())
		return "", err
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
	return string(toJsonData), nil
}
func selectingDataToBuildReport(completeQuery string) ([]string, error) {
	connection := getConnectionDB()
	query, err := connection.Query(completeQuery)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	var informationForReport []string
	for query.Next() {
		var dataReport dataPatientHC
		err = query.Scan(&dataReport.IdPatient, &dataReport.TypeId, &dataReport.DateClinicHistory, &dataReport.ActualDateRegistry, &dataReport.PatientNames, &dataReport.PatientLastnames, &dataReport.HasError)
		if err != nil {
			fmt.Println("Error scanning : " + err.Error())
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
	// getting the number of registries actually
	query, err = connection.Query("SELECT COUNT(IdPatient) FROM " + DATABASE_IN_USE)
	if err != nil {
		fmt.Println("Error in query getting amount of: " + err.Error())
	}
	var amount dataPatientHC
	for query.Next() {
		err = query.Scan(&amount.IdPatient)
		if err != nil {
			fmt.Println("Error scanning the data in getting amount of: " + err.Error())
		}
	}
	toJsonAmount, err := json.Marshal(amount)
	if err != nil {
		fmt.Println("Error marshallig the data: " + err.Error())
	}
	informationForReport = append(informationForReport, string(toJsonAmount))
	defer connection.Close()
	return informationForReport, nil
}
func getReport(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		dateStart := r.URL.Query().Get("date-start")
		dateEnd := r.URL.Query().Get("date-end")
		checkPatientErrors := r.URL.Query().Get("check-only-p-errors")
		generateReportBy := r.URL.Query().Get("gen-by")

		query := fmt.Sprintf("CALL GET_REPORT(%s, %s, %s, %s)", generateReportBy, "'"+dateStart+"'", "'"+dateEnd+"'", checkPatientErrors)
		dataRequest, err := selectingDataToBuildReport(query)
		if err != nil {
			fmt.Println("Error: " + err.Error())
		}
		fmt.Fprint(w, convertDataInJson(dataRequest))
	}
}
func getReportByPatient(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		argumentsQuery := r.URL.Query().Get("query-string")
		argumentsField := r.URL.Query().Get("query-field")

		query := fmt.Sprintf("CALL INFO_BY_PATIENT(%s, %s)", argumentsQuery, argumentsField)
		dataRequest, err := selectingDataToBuildReport(query)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Fprint(w, convertDataInJson(dataRequest))
	}
}
func convertDataInJson(dataRequest []string) string {
	toJson, err := json.Marshal(dataRequest)
	if err != nil {
		fmt.Println("Error marshalling: " + err.Error())
	}
	return string(toJson)
}
func reportInExcel(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var dataExcel setDataExcel
		contentRequest, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error: " + err.Error())
		}
		json.Unmarshal(contentRequest, &dataExcel)
		pathDownload, err := createExcelReport(dataExcel)
		if err != nil {
			fmt.Println("Error " + err.Error())
		}
		linkDownload := struct{ Link string }{Link: strings.Replace(pathDownload, "..", "", 1)}
		if ld, err := json.Marshal(linkDownload); err != nil {
			fmt.Println("Error marshalling link download: " + err.Error())
		} else {
			fmt.Fprint(w, string(ld))
		}
	}
}
func createExcelReport(contentData setDataExcel) (string, error) {
	contentHeaders := []string{"Documento N°", "Tipo ID", "Fecha de historia", "Fecha de registro", "Nombres", "Apellidos", "¿Paciente con error?"}
	indexColumns := []string{"A", "B", "C", "D", "E", "F", "G"}
	columnsWidth := []float64{21.86, 10.71, 23.86, 23.86, 27, 27, 26}
	sheetName := "Reporte_Pacientes"
	reportInExcel := excelize.NewFile()
	reportInExcel.NewSheet(sheetName)
	reportInExcel.DeleteSheet("Sheet1")
	// set headers
	for i := 0; i < len(contentHeaders); i++ {
		reportInExcel.SetCellValue(sheetName, indexColumns[i]+strconv.Itoa(1), contentHeaders[i])
		reportInExcel.SetColWidth(sheetName, indexColumns[i], indexColumns[i], columnsWidth[i])
	}
	// creating style
	styleReport, err := reportInExcel.NewStyle(`
		{
			"alignment":{"horizontal":"center"}, 
			"font":{"bold":true, "color":"#fffff"}, 
			"fill":{"type":"pattern", "color":["#00ADEF"], "pattern":1}, 
			"border":[
				{"type":"left", "color":"#000000", "style":1},
				{"type":"right", "color":"#000000", "style":1},
				{"type":"top", "color":"#000000", "style":1},
				{"type":"bottom", "color":"#000000", "style":1}]
		}`)
	if err != nil {
		fmt.Println("Error applying styles: " + err.Error())
	}
	// applying styles
	reportInExcel.SetCellStyle(sheetName, indexColumns[0]+strconv.Itoa(1), indexColumns[len(indexColumns)-1]+strconv.Itoa(1), styleReport)
	// inserting data
	for i := 0; i < len(contentData.DataExcel); i++ {
		var dataPatientExcel dataPatientHC
		err := json.Unmarshal([]byte(contentData.DataExcel[i]), &dataPatientExcel)
		if err != nil {
			return "", err
		}
		for j := 0; j < len(contentHeaders); j++ {
			var contentString string
			switch j {
			case 0:
				contentString = strconv.Itoa(dataPatientExcel.IdPatient)
			case 1:
				contentString = dataPatientExcel.TypeId
			case 2:
				contentString = dataPatientExcel.DateClinicHistory
			case 3:
				contentString = dataPatientExcel.ActualDateRegistry
			case 4:
				contentString = dataPatientExcel.PatientNames
			case 5:
				contentString = dataPatientExcel.PatientLastnames
			default:
				contentString = dataPatientExcel.HasError
			}
			reportInExcel.SetCellValue(sheetName, indexColumns[j]+strconv.Itoa(i+2), contentString)
		}
	}

	err = reportInExcel.SaveAs("../public/reports/Reporte.xlsx")
	if err != nil {
		return "", nil
	}
	return reportInExcel.Path, err
}
func main() {
	publicElementsApp := http.FileServer(http.Dir("../public"))
	http.Handle("/public/", http.StripPrefix("/public/", publicElementsApp))
	fmt.Println("Using the database: " + DATABASE_IN_USE)
	http.HandleFunc("/record-patient", setPatientRecord)
	http.HandleFunc("/get-data-patient", patientHosvital)
	http.HandleFunc("/get-information-from-patient", getReport)
	http.HandleFunc("/get-information-by-patient", getReportByPatient)
	http.HandleFunc("/get-report-in-excel", reportInExcel)
	http.HandleFunc("/Yuls", app)
	http.ListenAndServe(":8005", nil)
}
