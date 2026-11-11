package simpleRoutes

import (
	"baize/app/simple/controller"

	"github.com/gin-gonic/gin"
)

func InitSimpleRouter(router *gin.RouterGroup) {
	simpleGroup := router.Group("/simple")
	simpleGroup.GET("/customer/options", controller.SimpleCustomerOptions)
}
