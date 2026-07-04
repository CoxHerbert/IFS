package agentRoutes

import (
	"baize/app/agent/controller"
	"baize/app/common/middlewares"

	"github.com/gin-gonic/gin"
)

func InitPortalShipmentAnalyzeRouter(router *gin.RouterGroup) {
	group := router.Group("/api/shipment")
	group.POST("/analyze", controller.AnalyzeShipment)
}

func InitShipmentAnalyzeManageRouter(router *gin.RouterGroup) {
	group := router.Group("/agent/shipment")
	group.POST("/analyze", middlewares.HasPermission("ifs:agent:chat"), controller.AnalyzeShipment)
}
