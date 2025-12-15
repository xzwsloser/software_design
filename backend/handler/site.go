package handler

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/xzwsloser/software_design/backend/dto"
	"github.com/xzwsloser/software_design/backend/service"
	"github.com/xzwsloser/software_design/backend/utils"
)

type SiteHandler struct {

}

var (
	siteService  = new(service.SiteService)
	siteHandler  = new(SiteHandler)
)

// @Description: 景点分页查询接口
// url: /query/list
// method: POST
func (*SiteHandler) SitePageQuery(c *gin.Context) {
	var pageQueryParam dto.ScrollRequest
	err := c.ShouldBindJSON(&pageQueryParam)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to bind request params"))
		return
	}

	site_infos, err := siteService.QueryByPageParams(&pageQueryParam)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to query info"))
		return
	}

	result := dto.ScrollResp[dto.SiteBasicInfo]{
		Data: site_infos,
		Total: int32(len(site_infos)),
	}

	c.JSON(http.StatusOK, dto.OkWithData(result))
}

// @Description: 查询指定景点
// url: /query/:siteIndex
// method: GET
func (*SiteHandler) SiteQueryByIndex(c *gin.Context) {
	siteIndexStr := c.Param("siteIndex")
	siteIndex, err := strconv.Atoi(siteIndexStr)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("site index not a number"))
		return
	}

	site, err := siteService.QueryByIndex(int32(siteIndex))
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("cannot find site"))
		return
	}

	c.JSON(http.StatusOK, dto.OkWithData(site))
}






