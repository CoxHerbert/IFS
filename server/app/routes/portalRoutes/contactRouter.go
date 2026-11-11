package portalRoutes

import (
	"baize/app/common/middlewares"
	"baize/app/portal/controller"

	"github.com/gin-gonic/gin"
)

func InitContactRouter(router *gin.RouterGroup) {
	group := router.Group("/portal")
	group.GET("/version", controller.GetVersion)
	group.POST("/contact", controller.SubmitContact)
}

func InitContactManageRouter(router *gin.RouterGroup) {
	group := router.Group("/portal/contact")
	group.GET("/list", middlewares.HasPermission("portal:contact:list"), controller.ContactList)
	group.GET("/export", middlewares.HasPermission("portal:contact:export"), controller.ContactExport)
	group.GET("/:contactId", middlewares.HasPermission("portal:contact:query"), controller.ContactGetInfo)
	group.PUT("", middlewares.HasPermission("portal:contact:edit"), controller.ContactEdit)
	group.DELETE("/:contactIds", middlewares.HasPermission("portal:contact:remove"), controller.ContactRemove)
}
