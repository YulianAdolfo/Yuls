package readerparams

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type sqlConnection struct {
	Server   string
	Port     string
	Username string
	Password string
	Database string
}

func initReader() (*os.File, error) {
	file, err := os.Open("../PARAMETERS/APP_PARAMETERS.json")
	return file, err
}
func ReadConnectionSqlParameters() (string, string, string, string, string) {
	file, err := initReader()
	if err != nil {
		fmt.Println("¡Error al leer parametros! ", err)
	}
	data := readParameters(file)
	server := data["SqlHosvitalConnection"].(map[string]interface{})["Server"].(string)
	username := data["SqlHosvitalConnection"].(map[string]interface{})["Username"].(string)
	password := data["SqlHosvitalConnection"].(map[string]interface{})["Password"].(string)
	port := data["SqlHosvitalConnection"].(map[string]interface{})["Port"].(string)
	database := data["SqlHosvitalConnection"].(map[string]interface{})["Database"].(string)
	defer file.Close()
	return server, username, password, port, database

}
func ReadConnectionMySqlParameters() (string, string, string, string, string, string) {
	file, err := initReader()
	if err != nil {
		fmt.Println("¡Error al leer parametros! ", err)
	}
	data := readParameters(file)
	username := data["MysqlRemoteConnection"].(map[string]interface{})["Username"].(string)
	password := data["MysqlRemoteConnection"].(map[string]interface{})["Password"].(string)
	typeConnection := data["MysqlRemoteConnection"].(map[string]interface{})["TypeConn"].(string)
	server := data["MysqlRemoteConnection"].(map[string]interface{})["Server"].(string)
	port := data["MysqlRemoteConnection"].(map[string]interface{})["Port"].(string)
	database := data["MysqlRemoteConnection"].(map[string]interface{})["Database"].(string)
	defer file.Close()
	return username, password, typeConnection, server, port, database

}
func readParameters(file *os.File) map[string]interface{} {
	valuesinJson, _ := io.ReadAll(file)
	var valuesFromJson map[string]interface{}
	json.Unmarshal([]byte(valuesinJson), &valuesFromJson)
	return valuesFromJson
}
