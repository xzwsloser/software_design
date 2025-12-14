package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xzwsloser/software_design/backend/handler"
	"github.com/xzwsloser/software_design/backend/utils"
	"fmt"
)

func main() {
	utils.LoadConfig("config.json")
	utils.InitLogger()

	r := gin.Default()
	handler.InitRouter(r)

	serverAddr := fmt.Sprintf("localhost:%d", utils.GetServerConfig().Port)
	utils.GetLogger().Infof("Gin Server Listen At %s", serverAddr)
	r.Run(serverAddr)
}