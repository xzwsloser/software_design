package utils

import "testing"

func TestLoadConfig(t *testing.T) {
	LoadConfig("../config.json")
	t.Log("-----server config-----")
	t.Log(serverConfig.Port)
	t.Log("-----database config-----")
	t.Log(databaseConfig.Addr)
	t.Log(databaseConfig.Port)
	t.Log(databaseConfig.User)
	t.Log(databaseConfig.Password)
	t.Log(databaseConfig.Database)
}