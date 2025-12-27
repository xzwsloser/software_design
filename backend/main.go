package main

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/xzwsloser/software_design/backend/cache"
	"github.com/xzwsloser/software_design/backend/handler"
	"github.com/xzwsloser/software_design/backend/model"
	"github.com/xzwsloser/software_design/backend/rpc"
	"github.com/xzwsloser/software_design/backend/service"
	"github.com/xzwsloser/software_design/backend/utils"
)

func main() {
	// 加载配置文件, 初始化日志, 初始化 oss client
	utils.LoadConfig("config.json")
	utils.InitLogger()
	utils.InitOssClient()

	// 初始化数据库
	model.InitMysqlClient()
	cache.InitRedisClient()

	// 连接到 grpc server
	rpc.NewGrpcClient(utils.GetGrpcConfig().Addr,
					  utils.GetGrpcConfig().Port)

	// 初始化 后端 -> 推荐系统 pipeline
	service.NewRecSysPipeline()

	// 初始化路由
	r := gin.Default()

	// 更换 gin 框架日志中间件
	r.Use(gin.LoggerWithWriter(io.MultiWriter(utils.GetLogger().Writer())))
	r.Use(gin.Recovery())

	handler.InitRouter(r)

	serverAddr := fmt.Sprintf("localhost:%d", utils.GetServerConfig().Port)
	utils.GetLogger().Infof("Gin Server Listen At %s", serverAddr)
	r.Run(serverAddr)
}


