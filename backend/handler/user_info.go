package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xzwsloser/software_design/backend/dto"
	"github.com/xzwsloser/software_design/backend/middleware"
	"github.com/xzwsloser/software_design/backend/model"
	"github.com/xzwsloser/software_design/backend/utils"
)

type UserInfoHandler struct {

}

var (
	userInfoHandler = new(UserInfoHandler)
)

// @Description: 获取用户信息(之后可以做标签选择)
// url: /userInfo/user
// method: GET
func (*UserInfoHandler) GetUserInfo(c *gin.Context) {
	userInfo, err := middleware.GetBasicUserInfo(c)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to get user info from JWT"))
		return
	}

	user, err := userService.GetCurrentUserInfo(userInfo.Username)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to get user info"))
		return
	}

	c.JSON(http.StatusOK, dto.OkWithData(user))
}

// @Description: 更新用户信息
// url: /userInfo/update
// method: POST
func (*UserInfoHandler) UpdateUserInfo(c *gin.Context) {
	var u model.User
	err := c.ShouldBindJSON(&u)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to parse request"))
		return
	}

	userInfo, err := middleware.GetBasicUserInfo(c)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to get user info from jwt token"))
		return
	}

	u.Username = userInfo.Username
	u.Id	   = userInfo.Id

	err = userService.UpdateUserInfo(&u)

	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to update userInfo"))
		return
	}

	c.JSON(http.StatusOK, dto.Ok())
}


