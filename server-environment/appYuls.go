package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
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
	HasError           bool
	IDPTN              int
	DESCRIPTION_ERROR  string
	DATE               string
}
type returnMessage struct {
	ContenMessage string
}
type sqlColumnsName struct {
	MPNom1, MPApe1, MPTDoc, MPCedu string
}
type resultPatientSqlServer struct {
	Names, Lastnames, TypId, DocId string
}
type setDataExcel struct {
	DataExcel []string
}

var DATABASE_IN_USE string
var BACKUP_FILE_NAME = getFilenameBackups() + getDate() + getExt()

const VERSION = "1.1.0"

// insert new patients
func backup(dataPatienStruct dataPatientHC) {
	line := dataPatienStruct.ActualDateRegistry + ";" + dataPatienStruct.DateClinicHistory + ";" + strconv.Itoa(dataPatienStruct.IdPatient) + ";" + dataPatienStruct.PatientNames + ";" + dataPatienStruct.PatientLastnames + ";" + dataPatienStruct.TypeId + ";" + strconv.FormatBool(dataPatienStruct.HasError) + ";" + dataPatienStruct.DESCRIPTION_ERROR + ";" + dataPatienStruct.DATE + "\n"
	if err := saveDataInLocalBackup(line); err != nil {
		log.Print("Error saving the data in local way: " + err.Error())
	}
}
func newClinicHistory(dataPatienStruct dataPatientHC) error {
	go backup(dataPatienStruct)
	connection := getConnectionDB()
	knowExistancePatient := thisPatientExists(strconv.Itoa(dataPatienStruct.IdPatient))
	if knowExistancePatient != 1 {
		insertQuery := fmt.Sprintf("INSERT INTO %s (actualDateRegistry, dateClinicHistory, IdPatient, patientNames, patientLastnames, typeId, hasError) VALUES (?,?,?,?,?,?,?)", DATABASE_IN_USE)
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
		// verify if the patient has some errors
		if dataPatienStruct.HasError {
			err = insertDigitErrors(dataPatienStruct.IdPatient, dataPatienStruct.DESCRIPTION_ERROR, dataPatienStruct.DateClinicHistory, connection)
			if err != nil {
				fmt.Println("Error: " + err.Error())
				return err
			}
		}
	} else {
		if dataPatienStruct.HasError {
			err := insertDigitErrors(dataPatienStruct.IdPatient, dataPatienStruct.DESCRIPTION_ERROR, dataPatienStruct.DateClinicHistory, connection)
			if err != nil {
				fmt.Println("Error: " + err.Error())
				return err
			}
		} else {
			return errors.New("already registered")
		}
	}
	defer connection.Close()
	return nil
}
func insertDigitErrors(id int, description, date string, connection *sql.DB) error {
	query := "INSERT INTO TABLE_ERRORS (IDPTN, DESCRIPTION_ERROR, DATE) VALUES (?,?, ?)"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	statement, err := connection.PrepareContext(ctx, query)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return err
	}
	_, err = statement.ExecContext(ctx, id, description, date)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return err
	}
	return nil
}
func thisPatientExists(id string) int {
	connectionToDatabase := getConnectionDB()
	query := fmt.Sprintf("SELECT EXISTS (SELECT * FROM "+DATABASE_IN_USE+" WHERE IdPatient = '%s') AS patientExist", id)
	var patientExist string
	err := connectionToDatabase.QueryRow(query).Scan(&patientExist)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	state, err := strconv.Atoi(patientExist)
	if err != nil {
		fmt.Println("Cannot convert to int: " + err.Error())
	}
	return state
}
func setPatientRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		bodyRequest, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error reading the body: " + err.Error())
		}
		var dataPatientHc dataPatientHC
		json.Unmarshal(bodyRequest, &dataPatientHc)
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
	developerName := struct{ DevName string }{DevName: "Diseñado y desarrollado por Yulian Adolfo Rojas - Versión: " + VERSION}
	appTemplate := template.Must(template.ParseFiles("../client-environment/app.html"))
	appTemplate.Execute(w, developerName)
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
func getPatientsByNameFromHosvital(patientName string, connectionSqlServer *sql.DB) (string, error) {
	contextConnection := context.Background()
	// check if the connection is alive
	err := connectionSqlServer.PingContext(contextConnection)
	if err != nil {
		fmt.Println("Error in ping connection to Hosvital " + err.Error())
		return "", err

	}
	// if the connection is alive so create the sql qery
	sqlGetInfo := "SELECT TOP 50\n" +
		"RTRIM(MPCedu),\n" +
		"RTRIM(CONCAT(CONCAT(LEFT(MPNom1, 1), LOWER(RIGHT(RTRIM(MPNom1), LEN(MPNom1)-1))),' ',IIF (LEN(RTRIM(MPNom2))=0,'', CONCAT(LEFT(MPNom2, 1), LOWER(RIGHT(RTRIM(MPNom2), LEN(MPNom2)-1)))))),\n+" +
		"RTRIM(CONCAT(CONCAT(LEFT(MPApe1, 1), LOWER(RIGHT(RTRIM(MPApe1), LEN(MPApe1)-1))),' ',IIF  (LEN(RTRIM(MPApe2))=0,'', CONCAT(LEFT(MPApe2, 1), LOWER(RIGHT(RTRIM(MPApe2), LEN(MPApe2)-1))))))\n" +
		"FROM CAPBAS WHERE CONCAT(RTRIM(MPNom1),' ',RTRIM(MPNom2), ' ', RTRIM(MPApe1), ' ', RTRIM(MPApe2)) LIKE " + "'" + patientName + "%'"
	rows, err := connectionSqlServer.QueryContext(contextConnection, sqlGetInfo)
	if err != nil {
		fmt.Println("Error connection to Hosvital: " + err.Error())
		return "", err
	}
	defer rows.Close()

	var patientListByName []interface{}
	for rows.Next() {
		var dataSqlServer sqlColumnsName
		err = rows.Scan(&dataSqlServer.MPCedu, &dataSqlServer.MPNom1, &dataSqlServer.MPApe1)
		if err != nil {
			fmt.Println("Error scannig data from sql-server: " + err.Error())
			return "", err
		}
		// asigning values to then convert then into json
		patientListByName = append(patientListByName, dataSqlServer)
	}
	defer connectionSqlServer.Close()
	toJsonData, err := json.Marshal(patientListByName)
	if err != nil {
		fmt.Println("Error marshalling the json: " + err.Error())
		return "", err
	}
	return string(toJsonData), nil
}
func selectingDataToBuildReport(completeQuery string) ([]string, error) {
	connection := getConnectionDB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query, err := connection.QueryContext(ctx, completeQuery)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	var informationForReport []string
	for query.Next() {
		var dataReport dataPatientHC
		//SELECT typeId, IdPatient, dateClinicHistory, actualDateRegistry, patientNames, patientLastnames, IDPTN, DESCRIPTION_ERROR, DATE
		err = query.Scan(&dataReport.TypeId, &dataReport.IdPatient, &dataReport.DateClinicHistory, &dataReport.ActualDateRegistry, &dataReport.PatientNames, &dataReport.PatientLastnames, &dataReport.HasError, &dataReport.IDPTN, &dataReport.DESCRIPTION_ERROR, &dataReport.DATE)
		if err != nil {
			fmt.Println(err.Error())
		}
		content, err := json.Marshal(dataPatientHC{
			TypeId:             dataReport.TypeId,
			IdPatient:          dataReport.IdPatient,
			DateClinicHistory:  dataReport.DateClinicHistory,
			ActualDateRegistry: dataReport.ActualDateRegistry,
			PatientNames:       dataReport.PatientNames,
			PatientLastnames:   dataReport.PatientLastnames,
			HasError:           dataReport.HasError,
			IDPTN:              dataReport.IDPTN,
			DESCRIPTION_ERROR:  dataReport.DESCRIPTION_ERROR,
			DATE:               dataReport.DATE,
		})
		if err != nil {
			fmt.Println("Error marshalling data: " + err.Error())
		}
		informationForReport = append(informationForReport, string(content))
	}
	// getting the number of registries actually
	var amount dataPatientHC
	err = connection.QueryRow("SELECT COUNT(IdPatient) FROM " + DATABASE_IN_USE).Scan(&amount.IdPatient)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	informationForReport = append(informationForReport, strconv.Itoa(amount.IdPatient))
	defer connection.Close()
	return informationForReport, nil
}
func getReport(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		dateStart := r.URL.Query().Get("date-start")
		dateEnd := r.URL.Query().Get("date-end")
		checkPatientErrors := r.URL.Query().Get("check-only-p-errors")
		generateReportBy := r.URL.Query().Get("gen-by")

		checkErrors, err := strconv.Atoi(checkPatientErrors)
		if err != nil {
			fmt.Println("Fail to convert to int data")
		}
		genby, err := strconv.Atoi(generateReportBy)
		if err != nil {
			fmt.Println("Fail to convert to int data")
		}

		query := PrepareQueryForReport(checkErrors, genby, dateStart, dateEnd)

		dataRequest, err := selectingDataToBuildReport(query)
		if err != nil {
			fmt.Println("Error: " + err.Error())
		}
		fmt.Fprint(w, convertDataInJson(dataRequest))
	}
}
func PrepareQueryForReport(containsOnlyErrors, typeDate int, dateStart, dateEnd string) string {
	/* 	0 - it means the information will be generated by history clinic date,
	if not, will be generated by registry date
	*/
	tableAndField := DATABASE_IN_USE
	dbId := DATABASE_IN_USE + "." + "IdPatient"
	var query string
	if typeDate == 0 {
		tableAndField = tableAndField + "." + "dateClinicHistory"
	} else {
		tableAndField = tableAndField + "." + "actualDateRegistry"
	}

	if containsOnlyErrors != 1 {
		query = fmt.Sprintf("SELECT typeId, IdPatient, dateClinicHistory, actualDateRegistry, patientNames, patientLastnames, hasError, IDPTN, DESCRIPTION_ERROR, DATE FROM %s LEFT JOIN TABLE_ERRORS ON %s = TABLE_ERRORS.IDPTN WHERE %s BETWEEN '%s' AND '%s' ORDER BY %s ASC", DATABASE_IN_USE, dbId, tableAndField, dateStart, dateEnd, dbId)
	} else {
		query = fmt.Sprintf("SELECT typeId, IdPatient, dateClinicHistory, actualDateRegistry, patientNames, patientLastnames, hasError, IDPTN, DESCRIPTION_ERROR, DATE FROM %s INNER JOIN TABLE_ERRORS WHERE %s = TABLE_ERRORS.IDPTN AND %s BETWEEN '%s' AND '%s' ORDER BY %s", DATABASE_IN_USE, dbId, tableAndField, dateStart, dateEnd, dbId)
	}
	return query
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
func patientNameHosvital(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		patientName := r.URL.Query().Get("username-patient")
		if patientName != "" {
			requestPatientNames, err := getPatientsByNameFromHosvital(patientName, sqlServerGetConnection())
			if err != nil {
				fmt.Fprint(w, responseClientError(err))
			}
			fmt.Fprint(w, requestPatientNames)
		}
	}
}
func createExcelReport(contentData setDataExcel) (string, error) {
	contentHeaders := []string{"Documento N°", "Tipo ID", "Fecha de historia", "Fecha de registro", "Nombres", "Apellidos", "¿Error de digitación?", "Documento", "Descripción de error", "Fecha Error"}
	indexColumns := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	columnsWidth := []float64{21.86, 10.71, 23.86, 23.86, 27, 27, 26, 35, 30, 30}
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
	for i := 0; i < len(contentData.DataExcel)-1; i++ {
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
			case 6:
				if dataPatientExcel.HasError {
					contentString = "SI"
				} else {
					contentString = "NO"
				}
			case 7:
				contentString = strconv.Itoa(dataPatientExcel.IDPTN)
			case 8:
				contentString = dataPatientExcel.DESCRIPTION_ERROR
			case 9:
				contentString = dataPatientExcel.DATE
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
func readDBInUse() string {
	fmt.Println("Leyendo base de datos para uso")
	choicedDB, err := ioutil.ReadFile("../PARAMETERS/DB_PRODUCTION.txt")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	return string(choicedDB)
}
func getDate() string {
	return time.Now().Format("01-02-2006")
}
func getFilenameBackups() string {
	return "Yuls-Backup-"
}
func getExt() string {
	return ".txt"
}
func getPathBackup() string {
	return "./"
}
func saveDataInLocalBackup(data string) error {
	file, err := os.OpenFile(BACKUP_FILE_NAME, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Print("Error: " + err.Error())
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(data); err != nil {
		return err
	}
	return nil
}
func main() {
	// Eval file to backup
	_, err := os.Stat(getPathBackup() + "/" + BACKUP_FILE_NAME)
	if os.IsNotExist(err) {
		fmt.Println("Creating backup file on specified path...")
		_, err := os.Create(BACKUP_FILE_NAME)
		if err != nil {
			log.Print("Error creating the backup file: " + err.Error())
		}
		fmt.Println("File created")
	}
	fmt.Println("Leyendo parametros de conexión...")
	const PATH = "PARAMETERS/ADDRESS_IP_AND_PORT.txt"
	content, err := ioutil.ReadFile("../" + PATH)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		fmt.Println("Error el leer los parametros de conexión")
	}
	connectionParams := strings.TrimSpace(string(content))
	if connectionParams != "" {
		DATABASE_IN_USE = readDBInUse()
		addressAnPort := strings.Split(connectionParams, ":")
		publicElementsApp := http.FileServer(http.Dir("../public"))
		http.Handle("/public/", http.StripPrefix("/public/", publicElementsApp))
		fmt.Println("Usando la base de datos: " + DATABASE_IN_USE)
		http.HandleFunc("/record-patient", setPatientRecord)
		http.HandleFunc("/get-data-patient", patientHosvital)
		http.HandleFunc("/get-information-from-patient", getReport)
		http.HandleFunc("/get-information-by-patient", getReportByPatient)
		http.HandleFunc("/get-report-in-excel", reportInExcel)
		http.HandleFunc("/data-patient-from-hosvital", patientNameHosvital)
		http.HandleFunc("/Yuls", app)
		// opening the browers
		go func() {
			fmt.Println("Abriendo navegador/explorador...")
			<-time.After(100 * time.Millisecond)
			err := exec.Command("explorer", "http://"+addressAnPort[0]+":"+addressAnPort[1]+"/"+"Yuls").Run()
			if err != nil {
				fmt.Println("---------------------- Error --------------------")
				log.Print(err)
			}
		}()
		http.ListenAndServe(":"+addressAnPort[1], nil)
	} else {
		fmt.Println("SIN PARAMETROS DE CONEXIÓN")
	}
}
