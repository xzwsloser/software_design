package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/xzwsloser/software_design/backend/dto"
	"github.com/xzwsloser/software_design/backend/middleware"
	"github.com/xzwsloser/software_design/backend/utils"
)

func InitRouter(r *gin.Engine) {
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "pong")
	})

	// 解决跨域问题
	r.Use(middleware.CorsAllow())

	api := r.Group("/api")

	// 用户登录/注册
	userRouter := api.Group("/user")
	{
		userRouter.POST("/login", userHandler.Login)
		userRouter.POST("/register", userHandler.Register)
	}

	userInfoRouter := api.Group("/userInfo", middleware.JwtAuth())
	{
		userInfoRouter.GET("/user", userInfoHandler.GetUserInfo)
		userInfoRouter.POST("/update", userInfoHandler.UpdateUserInfo)
	}

	// 景点相关接口
	siteRouter := api.Group("/site", middleware.JwtAuth())
	{
		siteRouter.POST("/query/list", siteHandler.SitePageQuery)
		siteRouter.GET("/query/:siteIndex", siteHandler.SiteQueryByIndex)
		siteRouter.POST("/query/siteList", siteHandler.QueryBySiteIndexList)
	}

	// 评论相关接口
	commentRouter := api.Group("/comment", middleware.JwtAuth())
	{
		commentRouter.POST("/positive/:siteIndex", commentHandler.QueryPositiveCommentByPage)
		commentRouter.POST("/negative/:siteIndex", commentHandler.QueryNegativeCommentByPage)
		commentRouter.GET("/count/positive/:siteIndex", commentHandler.CountPositiveComment)
		commentRouter.GET("/count/negative/:siteIndex", commentHandler.CountNegativeComment)
	}

	// 点赞相关接口
	likeRouter := api.Group("/like", middleware.JwtAuth())
	{
		likeRouter.GET("/like/:siteIndex", likeHandler.Like)
		likeRouter.GET("/cancel/:siteIndex", likeHandler.CancelLike)
		likeRouter.GET("/isLike/:siteIndex", likeHandler.IsLikeSite)
		likeRouter.GET("/siteList", likeHandler.QuerySiteLikedByUser)
		likeRouter.GET("/userList/:siteIndex", likeHandler.QueryUserListLikedSite)
	}

	// 浏览记录相关接口
	viewRouter := api.Group("/view", middleware.JwtAuth())
	{
		viewRouter.GET("/view/:siteIndex", viewHandler.View)
		viewRouter.GET("/siteList", viewHandler.GetVisitedSiteList)
		viewRouter.GET("/userList/:siteIndex", viewHandler.GetUserList)
	}

	// OSS 图片获取接口
	ossRouter := api.Group("/oss", middleware.JwtAuth())
	{
		ossRouter.GET("/touristType", ossHandler.GetTotalTouristTypePic)
	}

	// 测试接口
	testRouter := api.Group("/test", middleware.JwtAuth())
	{
		// 测试 jwt 
		testRouter.GET("/jwt", func(c *gin.Context) {
			userInfo, err := middleware.GetBasicUserInfo(c)
			if err != nil {
				utils.GetLogger().Error(err.Error())
				c.JSON(http.StatusOK, dto.Fail("Failed to get user Info"))
				return
			}

			c.JSON(http.StatusOK, dto.OkWithData(userInfo))
		})
	}
}


