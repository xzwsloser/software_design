package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xzwsloser/software_design/backend/dto"
	"github.com/xzwsloser/software_design/backend/middleware"
	"github.com/xzwsloser/software_design/backend/service"
	"github.com/xzwsloser/software_design/backend/utils"
)

type RecSysHandler struct {

}

var (
	recSysService = new(service.RecSysService)
	recSysHandler = new(RecSysHandler)
)

// @Description: 查询推荐景点列表
// url: /rec/siteIdxList
// method: GET
func (*RecSysHandler) GetRecSiteIdxList(c *gin.Context) {
	userInfo, err := middleware.GetBasicUserInfo(c)

	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to get user info from jwt token"))
		return
	}

	userId := int(userInfo.Id)

	siteIdxList, err := recSysService.QueryRecommandSiteIdxs(userId)
	if err != nil {
		utils.GetLogger().Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail("Failed to get recommand site idx list"))
		return
	}

	resp := dto.ScrollResp[int]{
		Data: siteIdxList,
		Total: int32(len(siteIdxList)),
	}

	c.JSON(http.StatusOK, dto.OkWithData(resp))

}





