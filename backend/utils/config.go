package utils

import (
	"encoding/json"
	"io"
	"os"
)

type ServerConfig struct {
	Port int32 `json:"port"`
}

type DatabaseConfig struct {
	Addr 		string  `json:"addr"`
	Port 		int32 	`json:"port"`
	User 		string  `json:"user"`
	Password	string	`json:"password"`
	Database	string  `json:"database"` 
}

var (
	serverConfig 	*ServerConfig = new(ServerConfig)
	databaseConfig  *DatabaseConfig = new(DatabaseConfig) 
)

func LoadConfig(filePath string) {
	file, err := os.Open(filePath)

	if err != nil {
		panic("Failed to read config file")
	}

	defer file.Close()

	content, err := io.ReadAll(file)
	
	if err != nil {
		panic("Failed to read config file content")
	}

	var configDict map[string]interface{}
	err = json.Unmarshal(content, &configDict)
	if err != nil {
		panic("Failed to Unmarshal config")
	}

	serverJsonCode, err := json.Marshal(configDict["server"].(map[string]interface{}))
	if err != nil {
		panic("Failed to Fetch server config")
	}

	err = json.Unmarshal(serverJsonCode, serverConfig)
	if err != nil {
		panic("Failed to transform to server config")
	}

	databaseJsonCode, err := json.Marshal(configDict["database"].(map[string]interface{}))
	if err != nil {
		panic("Failed to Fetch database config")
	}

	err = json.Unmarshal(databaseJsonCode, databaseConfig)
	if err != nil {
		panic("Failed to transform to database config")
	}
}

func GetServerConfig() *ServerConfig {
	return serverConfig
}

func GetDatabaseConfig() *DatabaseConfig {
	return databaseConfig
}
