package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xzwsloser/software_design/backend/dto"
	"github.com/xzwsloser/software_design/backend/utils"
)

type OssHandler struct {

}

var (
	ossHandler	= new(OssHandler)
)

// @Description: 获取总共的旅游动机分布图 url
// url: /oss/touristType
// method: GET
func(*OssHandler) GetTotalTouristTypePic(c *gin.Context) {
	totalTouristTypeUrl := utils.GetTotalTouristTypePic()
	c.JSON(http.StatusOK, dto.OkWithData(totalTouristTypeUrl))
}
