package freightRoutes

import (
	"baize/app/common/middlewares"
	"baize/app/freight/controller"
	"github.com/gin-gonic/gin"
)

func InitReceiptManageRouter(router *gin.RouterGroup) {
	group:=router.Group("/freight/receipt")
	group.GET("/list",middlewares.HasPermission("freight:receipt:list"),controller.ReceiptList)
	group.GET("/:receiptId",middlewares.HasPermission("freight:receipt:query"),controller.ReceiptGet)
	group.POST("",middlewares.HasPermission("freight:receipt:add"),controller.ReceiptCreate)
	group.DELETE("/:receiptId",middlewares.HasPermission("freight:receipt:remove"),controller.ReceiptRemove)
}

func InitPaymentDeclarationManageRouter(router *gin.RouterGroup){
	group:=router.Group("/freight/payment-declaration")
	group.GET("/list",middlewares.HasPermission("freight:declaration:list"),controller.PaymentDeclarationList)
	group.GET("/:declarationId",middlewares.HasPermission("freight:declaration:query"),controller.PaymentDeclarationGet)
	group.POST("/:declarationId/approve",middlewares.HasPermission("freight:declaration:review"),controller.PaymentDeclarationApprove)
	group.POST("/:declarationId/reject",middlewares.HasPermission("freight:declaration:review"),controller.PaymentDeclarationReject)
}
