package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "pong")
	})
}