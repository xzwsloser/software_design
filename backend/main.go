package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xzwsloser/software_design/backend/cache"
	"github.com/xzwsloser/software_design/backend/handler"
	"github.com/xzwsloser/software_design/backend/model"
	"github.com/xzwsloser/software_design/backend/utils"
)

func main() {
	utils.LoadConfig("config.json")
	utils.InitLogger()
	utils.InitOssClient()

	model.InitMysqlClient()
	cache.InitRedisClient()

	r := gin.Default()
	handler.InitRouter(r)

	serverAddr := fmt.Sprintf("localhost:%d", utils.GetServerConfig().Port)
	utils.GetLogger().Infof("Gin Server Listen At %s", serverAddr)
	r.Run(serverAddr)
}
