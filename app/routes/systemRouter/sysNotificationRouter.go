package systemRouter

import (
	"baize/app/common/middlewares"
	"baize/app/notification/controller"

	"github.com/gin-gonic/gin"
)

func InitSysNotificationRouter(router *gin.RouterGroup) {
	group := router.Group("/system/notification")
	group.GET("/list", middlewares.HasPermission("system:notification:list"), controller.NotificationList)
	group.GET("/unread-count", controller.NotificationUnreadCount)
	group.PUT("/:notificationId/read", middlewares.HasPermission("system:notification:edit"), controller.NotificationRead)
	group.PUT("/read-all", middlewares.HasPermission("system:notification:edit"), controller.NotificationReadAll)
	group.DELETE("/:notificationIds", middlewares.HasPermission("system:notification:remove"), controller.NotificationRemove)
}
