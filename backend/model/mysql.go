package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "gorm.io/driver/mysql"
	"github.com/xzwsloser/software_design/backend/utils"
)

const (
	DATABASE = "mysql"
)

var (
	_mysqlClient *gorm.DB
)

func InitMysqlClient() {
	databaseConfig := utils.GetDatabaseConfig()
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
					   databaseConfig.User,
					   databaseConfig.Password,
					   databaseConfig.Addr,
					   databaseConfig.Port,
					   databaseConfig.Database)

	db, err := gorm.Open(DATABASE, url)

	if err != nil {
		utils.GetLogger().Error(err.Error())
		panic("Cannot connect to Mysql")
	}

	_mysqlClient = db
}

func GetMySqlClient() *gorm.DB {
	return _mysqlClient
}
