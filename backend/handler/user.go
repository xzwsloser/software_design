package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/xzwsloser/software_design/backend/dto"
	"github.com/xzwsloser/software_design/backend/model"
	"github.com/xzwsloser/software_design/backend/service"
	"github.com/xzwsloser/software_design/backend/utils"
)

type UserHandler struct {

}

var (
	userHandler *UserHandler         = new(UserHandler)
	userService *service.UserService = new(service.UserService)
)

// @Description: 用户登录接口, 接受用户信息 返回 jwt token
// POST
func (*UserHandler) Login(c *gin.Context) {
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to bind user info"))
		return 
	}

	jwtToken, err := userService.Login(&user)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK,dto.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OkWithData(jwtToken))
}


// @Description: 用户注册接口, 接受用户信息 返回 jwt token
// POST
func (*UserHandler) Register(c *gin.Context) {
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to bind user info"))
		return
	}

	jwtToken, err := userService.Register(&user)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OkWithData(jwtToken))
}


