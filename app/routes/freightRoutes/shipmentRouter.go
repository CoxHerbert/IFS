package freightRoutes

import (
	"baize/app/common/middlewares"
	"baize/app/freight/controller"

	"github.com/gin-gonic/gin"
)

func InitShipmentManageRouter(router *gin.RouterGroup) {
	group := router.Group("/freight/shipment")
	group.GET("/list", middlewares.HasPermission("freight:shipment:list"), controller.ShipmentList)
	group.POST("/import", middlewares.HasPermission("freight:shipment:import"), controller.ShipmentImport)
	group.GET("/:shipmentId", middlewares.HasPermission("freight:shipment:query"), controller.ShipmentGetInfo)
	group.PUT("/:shipmentId/status", middlewares.HasPermission("freight:shipment:edit"), controller.ShipmentUpdateStatus)
	group.POST("/:shipmentId/confirm", middlewares.HasPermission("freight:shipment:confirm"), controller.ShipmentConfirm)
	group.GET("/:shipmentId/share", middlewares.HasPermission("freight:shipment:share"), controller.ShipmentShareInfo)
	group.DELETE("/:shipmentIds", middlewares.HasPermission("freight:shipment:remove"), controller.ShipmentRemove)
}

func InitPortalShipmentRouter(router *gin.RouterGroup) {
	group := router.Group("/portal/shipment")
	group.GET("/share/:token", controller.PortalShipmentShare)
}
