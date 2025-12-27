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

type JwtConfig struct {
	SerectKey	string 	`json:"serect"`
	Issuer		string 	`json:"issuer"`
}

type RedisConfig struct {
	Addr 	string	`json:"addr"`
	Port 	int32	`json:"port"`
	DB 		int32	`json:"db"`
}

type OssConfig struct {
	Url 		string	`json:"url"`
	AppId 		string	`json:"appId"`
	Bucket		string	`json:"bucket"`
	SecretId	string	`json:"secretId"`
	SecretKey	string	`json:"secretKey"`
}

type GrpcConfig struct {
	Addr 	string	`json:"addr"`
	Port	int		`json:"port"`
}

var (
	serverConfig 	*ServerConfig 	= new(ServerConfig)
	databaseConfig  *DatabaseConfig = new(DatabaseConfig) 
	jwtConfig 		*JwtConfig 		= new(JwtConfig)
	redisConfig 	*RedisConfig    = new(RedisConfig)
	ossConfig       *OssConfig      = new(OssConfig)
	grpcConfig		*GrpcConfig		= new(GrpcConfig)
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

	jwtJsonCode, err := json.Marshal(configDict["jwt"].(map[string]interface{}))
	if err != nil {
		panic("Failed to Fetch Jwt Config")
	}

	err = json.Unmarshal(jwtJsonCode, jwtConfig)
	if err != nil {
		panic("Failed to transform to jwt config")
	}

	redisJsonCode, err := json.Marshal(configDict["redis"].(map[string]interface{}))
	if err != nil {
		panic("Failed to Fetch Redis Config")
	}

	err = json.Unmarshal(redisJsonCode, redisConfig)
	if err != nil {
		panic("Failed to transform to redis config")
	}

	ossJsonCode, err := json.Marshal(configDict["oss"].(map[string]interface{}))
	if err != nil {
		panic("Failed to Fetch oss Config")
	}

	err = json.Unmarshal(ossJsonCode, ossConfig)
	if err != nil {
		panic("Failed to transform to oss Config")
	}

	grpcJsonCode, err := json.Marshal(configDict["grpc"].(map[string]interface{}))
	if err != nil {
		panic("Failed to Fetch grpc Config")
	}

	err = json.Unmarshal(grpcJsonCode, grpcConfig)
	if err != nil {
		panic("Failed to transform to grpc Config")
	}

}

func GetServerConfig() *ServerConfig {
	return serverConfig
}

func GetDatabaseConfig() *DatabaseConfig {
	return databaseConfig
}

func GetJwtConfig() *JwtConfig {
	return jwtConfig
}

func GetRedisConfig() *RedisConfig {
	return redisConfig
}

func GetOssConfig() *OssConfig {
	return ossConfig
}

func GetGrpcConfig() *GrpcConfig {
	return grpcConfig
}

