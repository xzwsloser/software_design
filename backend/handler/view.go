package handler

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/xzwsloser/software_design/backend/dto"
	"github.com/xzwsloser/software_design/backend/middleware"
	"github.com/xzwsloser/software_design/backend/service"
	"github.com/xzwsloser/software_design/backend/utils"
)

type ViewHandler struct {

}

var (
	viewService = new(service.ViewService)
	viewHandler = new(ViewHandler)
)


// @Description: 浏览景点记录
// url: /view/view/:siteIndex
// method: GET
func (*ViewHandler) View(c *gin.Context) {
	siteIndexStr := c.Param("siteIndex")
	siteIndex, err := strconv.Atoi(siteIndexStr)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("site index not a number"))
		return
	}

	userInfo, err := middleware.GetBasicUserInfo(c)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("cannot get user info from JWT"))
		return
	}

	userId := int(userInfo.Id)

	result, err := viewService.View(userId, siteIndex)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to record view"))
		return
	}

	c.JSON(http.StatusOK, dto.OkWithData(result))
}

// @Description: 获取当前用户浏览列表
// url: /view/siteList/
// method: GET
func (*ViewHandler) GetVisitedSiteList(c *gin.Context) {
	userInfo, err := middleware.GetBasicUserInfo(c)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to get user info from JWT"))
		return
	}

	userId := int(userInfo.Id)

	siteList, err := viewService.QueryVisitedSiteList(userId)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to get site list"))
		return
	}

	resp := &dto.ScrollResp[int]{
		Data: siteList,
		Total: int32(len(siteList)),
	}

	c.JSON(http.StatusOK, dto.OkWithData(resp))
}

// @Description: 获取到浏览个某一个景点的全部用户列表
// url /view/userList/:siteIndex
// method: POST
func (*ViewHandler) GetUserList(c *gin.Context) {
	siteIndexStr := c.Param("siteIndex")
	siteIndex, err := strconv.Atoi(siteIndexStr)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("site index not a number"))
		return
	}

	results, err := viewService.QueryUserListed(siteIndex)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to query view user list"))
		return
	}

	resp := dto.ScrollResp[int] {
		Data: results,
		Total: int32(len(results)),
	}

	c.JSON(http.StatusOK, dto.OkWithData(resp))
}


