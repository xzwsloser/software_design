package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xzwsloser/software_design/backend/dto"
	"github.com/xzwsloser/software_design/backend/model"
	"github.com/xzwsloser/software_design/backend/service"
	"github.com/xzwsloser/software_design/backend/utils"
)

type CommentHandler struct {

}

var (
	commentService = new(service.CommentService)
	commentHandler = new(CommentHandler)
)

// @Description: 分页查询正面评论
// url: /comment/positive/:siteIndex
// method: POST
func (*CommentHandler) QueryPositiveCommentByPage(c *gin.Context) {
	siteIndexStr := c.Param("siteIndex")
	siteIndex, err := strconv.Atoi(siteIndexStr)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("site Index not a number"))
		return
	}

	var pageQueryParam dto.ScrollRequest
	err = c.ShouldBindJSON(&pageQueryParam)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("cannot bind query params"))
		return
	}

	comments, err := commentService.QueryPositiveCommentByPage(int32(siteIndex),
															   pageQueryParam.PageIndex,
															   pageQueryParam.PageSize)			

	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("query positive comment failed"))
		return 
	}

   resp := dto.ScrollResp[model.CommentPositive] {
	Data: comments,
	Total: int32(len(comments)),
   }

   c.JSON(http.StatusOK, dto.OkWithData(resp))
}


// @Description: 分页查询负面评论
// url: /comment/negative/:siteIndex
// method: POST
func (*CommentHandler) QueryNegativeCommentByPage(c *gin.Context) {
	siteIndexStr := c.Param("siteIndex")
	siteIndex, err := strconv.Atoi(siteIndexStr)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("site index not a number"))
		return
	}

	var pageQueryParam dto.ScrollRequest
	err = c.ShouldBindJSON(&pageQueryParam)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("cannot bind page query params"))
		return
	}

	comments, err := commentService.QueryNegativeCommentByPage(int32(siteIndex),
														 	   pageQueryParam.PageIndex,
															   pageQueryParam.PageSize)

    if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("query negative comment failed"))
		return
	}

	resp := dto.ScrollResp[model.CommentNegative]{
		Data: comments,
		Total: int32(len(comments)),
	}

	c.JSON(http.StatusOK, dto.OkWithData(resp))
}

// @Description: 统计好评数量
// url: /comment/count/positive/:siteIndex
// method: GET
func (*CommentHandler) CountPositiveComment(c *gin.Context) {
	siteIndexStr := c.Param("siteIndex")
	siteIndex, err := strconv.Atoi(siteIndexStr)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("siteIndex not a number"))
		return
	}

	count, err := commentService.CountPositiveComment(int32(siteIndex))
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("count positive comment failed!"))
		return
	}

	c.JSON(http.StatusOK, dto.OkWithData(count))
}

// @Description: 统计差评数量
// url: /comment/count/negative/:siteIndex
// method: GET
func (*CommentHandler) CountNegativeComment(c *gin.Context) {
	siteIndexStr := c.Param("siteIndex")
	siteIndex, err := strconv.Atoi(siteIndexStr)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("siteIndex not a number"))
		return
	}

	count, err := commentService.CountNegativeComment(int32(siteIndex))
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("count negative comment failed!"))
		return
	}

	c.JSON(http.StatusOK, dto.OkWithData(count))
}
