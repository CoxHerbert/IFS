package cmsRoutes

import (
	"baize/app/cms/controller"
	"baize/app/common/middlewares"

	"github.com/gin-gonic/gin"
)

func InitArticleManageRouter(router *gin.RouterGroup) {
	group := router.Group("/cms/article")
	group.GET("/list", middlewares.HasPermission("cms:article:list"), controller.ArticleList)
	group.GET("/:articleId", middlewares.HasPermission("cms:article:query"), controller.ArticleGetInfo)
	group.POST("", middlewares.HasPermission("cms:article:add"), controller.ArticleAdd)
	group.PUT("", middlewares.HasPermission("cms:article:edit"), controller.ArticleEdit)
	group.DELETE("/:articleIds", middlewares.HasPermission("cms:article:remove"), controller.ArticleRemove)
	group.POST("/upload-image", middlewares.HasPermission("cms:article:list"), controller.ArticleUploadImage)
}
