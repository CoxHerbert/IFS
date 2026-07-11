package portalRoutes

import (
	"baize/app/cms/controller"

	"github.com/gin-gonic/gin"
)

func InitArticleRouter(router *gin.RouterGroup) {
	group := router.Group("/portal/articles")
	group.GET("", controller.PortalArticleList)
	group.GET("/:slug", controller.PortalArticleDetail)
}
