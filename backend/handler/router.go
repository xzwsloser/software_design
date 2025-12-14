package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/xzwsloser/software_design/backend/dto"
	"github.com/xzwsloser/software_design/backend/middleware"
	"github.com/xzwsloser/software_design/backend/utils"
)

func InitRouter(r *gin.Engine) {
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "pong")
	})

	// 用户登录/注册
	userRouter := r.Group("/user")
	{
		userRouter.POST("/login", userHandler.Login)
		userRouter.POST("/register", userHandler.Register)
	}

	// 测试接口
	testRouter := r.Group("/test", middleware.JwtAuth())
	{
		// 测试 jwt 
		testRouter.GET("/jwt", func(c *gin.Context) {
			userInfo, err := middleware.GetBasicUserInfo(c)
			if err != nil {
				utils.GetLogger().Error(err.Error())
				c.JSON(http.StatusOK, dto.Fail("Failed to get user Info"))
				return
			}

			c.JSON(http.StatusOK, dto.OkWithData(userInfo))
		})
	}
}
