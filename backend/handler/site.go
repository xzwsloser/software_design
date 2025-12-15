package handler

import (
	"net/http"

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

