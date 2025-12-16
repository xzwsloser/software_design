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

type LikeHandler struct {

}

var (
	likeService = new(service.LikeService)
	likeHandler = new(LikeHandler)
)

// @Decription: 点赞接口
// url: /like/like/:siteIndex
// method: GET
func (*LikeHandler) Like(c *gin.Context) {
	siteIndexStr := c.Param("siteIndex")
	siteIndex, err := strconv.Atoi(siteIndexStr)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Site Index not a number"))
		return
	}

	userInfo, err := middleware.GetBasicUserInfo(c)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to Get User Info From JWT"))
		return
	}

	userId := int(userInfo.Id)

	err = likeService.Like(userId, siteIndex)

	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to record like connection"))
		return
	}

	c.JSON(http.StatusOK, dto.Ok())
}

// @Description: 取消点赞接口
// url: /like/cancel/:siteIndex
// method: GET
func (*LikeHandler) CancelLike(c *gin.Context) {
	siteIndexStr := c.Param("siteIndex")
	siteIndex, err := strconv.Atoi(siteIndexStr)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Site Index not a number"))
		return
	}

	userInfo, err := middleware.GetBasicUserInfo(c)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to Get User Info From JWT"))
		return
	}

	userId := int(userInfo.Id)

	err = likeService.CancelLike(userId, siteIndex)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to cancel like connection"))
		return
	}

	c.JSON(http.StatusOK, dto.Ok())
}

// @Description: 查询用户是否点赞某一个景点
// url: /like/isLike/:siteIndex
// method: GET
func (*LikeHandler) IsLikeSite(c *gin.Context) {
	siteIndexStr := c.Param("siteIndex")
	siteIndex, err := strconv.Atoi(siteIndexStr)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Site Index not a number"))
		return
	}

	userInfo, err := middleware.GetBasicUserInfo(c)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to Get User Info From JWT"))
		return
	}

	userId := int(userInfo.Id)

	result, err := likeService.QueryIsLikedByUser(userId, siteIndex)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to query is liked"))
		return
	}

	c.JSON(http.StatusOK, dto.OkWithData(result))
}

// @Description: 查询用户点赞的景点列表
// url: /like/siteList
// method: GET
func (*LikeHandler) QuerySiteLikedByUser(c *gin.Context) {
	userInfo, err := middleware.GetBasicUserInfo(c)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Cannot get user id from JWT"))
		return
	}

	userId := int(userInfo.Id)

	results, err := likeService.QueryLikeOfUser(userId)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to query liked site list"))
		return
	}

	result := dto.ScrollResp[int]{
		Data: results,
		Total: int32(len(results)),
	}

	c.JSON(http.StatusOK, dto.OkWithData(result))
}

// @Description: 喜欢某一个景点的用户列表
// url: /like/userList/:siteIndex
// method: GET
func (*LikeHandler) QueryUserListLikedSite(c *gin.Context) {
	siteIndexStr := c.Param("siteIndex")
	siteIndex, err := strconv.Atoi(siteIndexStr)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("siteIndex not a number"))
		return
	}

	results, err := likeService.QueryLikeOfSite(siteIndex)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to query user list"))
		return 
	}

	result := dto.ScrollResp[int] {
		Data: results,
		Total: int32(len(results)),
	}

	c.JSON(http.StatusOK, dto.OkWithData(result))
}
