package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Configuration struct {
	Host              string `json:"host"`
	Port              string `json:"port"`
	Database          string `json:"database"`
	URL               string `json:"url"`
	ConfigurationType string `json:"configurationType"`
	Type              string `json:"type"`
	AuthModel         string `json:"auth-model"`
	User              string `json:"user"`
	Password          string `json:"password"`
}

type Connection struct {
	Provider      string        `json:"provider"`
	Driver        string        `json:"driver"`
	Name          string        `json:"name"`
	SavePassword  bool          `json:"save-password"`
	Configuration Configuration `json:"configuration"`
}

type ConnectionType struct {
	Name                    string `json:"name"`
	Color                   string `json:"color"`
	Description             string `json:"description"`
	AutoCommit              bool   `json:"auto-commit"`
	ConfirmExecute          bool   `json:"confirm-execute"`
	ConfirmDataChange       bool   `json:"confirm-data-change"`
	SmartCommit             bool   `json:"smart-commit"`
	SmartCommitRecover      bool   `json:"smart-commit-recover"`
	AutoCloseTransactions   bool   `json:"auto-close-transactions"`
	CloseTransactionsPeriod int    `json:"close-transactions-period"`
}

type Config struct {
	Folders         map[string]interface{}    `json:"folders"`
	Connections     map[string]Connection     `json:"connections"`
	ConnectionTypes map[string]ConnectionType `json:"connection-types"`
}

func Setup_dbeaver(credentials Credentials, dbname string) {
	fmt.Println("Setting up DBeaver")
	path := os.Getenv("DBEAVER_CONFIG_PATH") + "/data-sources.json"

	//retrieve the dbeaver connection file
	jsonFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer jsonFile.Close()
	fmt.Printf("file: %vn", jsonFile)

	fmt.Println("File opened successfully", jsonFile.Name())

	byteValue, _ := io.ReadAll(jsonFile)

	// we initialize our Users array
	var dataSources Config

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &dataSources)

	// fmt.Println("DataSources: ", dataSources)

	connection := Connection{
		Provider:     "postgresql",
		Driver:       "postgres-jdbc",
		Name:         credentials.Target.Name,
		SavePassword: true,
		Configuration: Configuration{
			Host:              "localhost",
			Port:              credentials.port,
			Database:          dbname,
			URL:               "jdbc:postgresql://localhost:" + credentials.port + "/" + dbname,
			ConfigurationType: "MANUAL",
			Type:              "boundary-cli",
			AuthModel:         "native",
			User:              credentials.username,
			Password:          credentials.password,
		},
	}
	connectionType := ConnectionType{
		Name:                    "boundary-cli",
		Color:                   "255,255,255",
		Description:             "Boundary CLI",
		AutoCommit:              true,
		ConfirmExecute:          false,
		ConfirmDataChange:       false,
		SmartCommit:             false,
		SmartCommitRecover:      false,
		AutoCloseTransactions:   true,
		CloseTransactionsPeriod: 1800,
	}

	dataSources.Connections[credentials.Target.ID] = connection
	// I want ConnectionTypes to be an empty map if it doesn't exist
	if dataSources.ConnectionTypes == nil {
		dataSources.ConnectionTypes = make(map[string]ConnectionType)
	}

	dataSources.ConnectionTypes["boundary-cli"] = connectionType

	fmt.Println("DataSources: ", dataSources)

	// convert the dataSources to json
	jsonData, err := json.Marshal(dataSources)
	if err != nil {
		fmt.Println("Error marshalling dataSources: ", err)
		return
	}
	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file: ", err)
		return
	}

	os.Exit(0)
}
