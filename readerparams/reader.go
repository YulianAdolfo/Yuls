package readerparams

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
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
func ReadDataInUsage() string {
	file, err := initReader()
	if err != nil {
		fmt.Println("¡Error al leer parametros! ", err)
	}
	data := readParameters(file)
	database := data["MysqlRemoteConnection"].(map[string]interface{})["Database"].(string)
	defer file.Close()
	return database
}
func ReadLocalNetwork() (string, string) {
	file, err := initReader()
	if err != nil {
		fmt.Println("¡Error al leer parametros! ", err)
	}
	data := readParameters(file)
	IP := data["LocalNetwork"].(map[string]interface{})["IP"].(string)
	Port := data["LocalNetwork"].(map[string]interface{})["Port"].(string)
	// testing if an IP is already set up
	// if so, the system will take the local IP, if not, it will take the set up IP
	if IP == "" || len(IP) <= 0 {
		IP = getLocalIP()
	}
	return IP, Port
}
func readParameters(file *os.File) map[string]interface{} {
	valuesinJson, _ := io.ReadAll(file)
	var valuesFromJson map[string]interface{}
	json.Unmarshal([]byte(valuesinJson), &valuesFromJson)
	return valuesFromJson
}
func getLocalIP() string {
	connection, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	localAddress := connection.LocalAddr().(*net.UDPAddr)
	return localAddress.IP.String()
}
