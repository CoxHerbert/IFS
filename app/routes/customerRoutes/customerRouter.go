package customerRoutes

import (
	"baize/app/common/middlewares"
	"baize/app/customer/controller"
	customermiddleware "baize/app/customer/middleware"

	"github.com/gin-gonic/gin"
)

func InitCustomerManageRouter(router *gin.RouterGroup) {
	group := router.Group("/customer")
	group.GET("/list", middlewares.HasPermission("customer:customer:list"), controller.CustomerList)
	group.GET("/options", middlewares.HasPermission("customer:account:list"), controller.CustomerOptions)
	group.POST("", middlewares.HasPermission("customer:customer:add"), controller.CustomerAdd)
	group.GET("/:customerId", middlewares.HasPermission("customer:customer:query"), controller.CustomerGetInfo)
	group.PUT("", middlewares.HasPermission("customer:customer:edit"), controller.CustomerEdit)
	group.DELETE("/:customerIds", middlewares.HasPermission("customer:customer:remove"), controller.CustomerRemove)
	group.GET("/:customerId/contact/list", middlewares.HasPermission("customer:customer:query"), controller.ContactList)

	contactGroup := router.Group("/customer/contact")
	contactGroup.POST("", middlewares.HasPermission("customer:customer:edit"), controller.ContactAdd)
	contactGroup.GET("/:contactId", middlewares.HasPermission("customer:customer:query"), controller.ContactGetInfo)
	contactGroup.PUT("", middlewares.HasPermission("customer:customer:edit"), controller.ContactEdit)
	contactGroup.DELETE("/:contactIds", middlewares.HasPermission("customer:customer:remove"), controller.ContactRemove)

	accountGroup := router.Group("/customer/account")
	accountGroup.GET("/list", middlewares.HasPermission("customer:account:list"), controller.AccountList)
	accountGroup.POST("", middlewares.HasPermission("customer:account:add"), controller.AccountAdd)
	accountGroup.GET("/:accountId", middlewares.HasPermission("customer:account:query"), controller.AccountGetInfo)
	accountGroup.PUT("", middlewares.HasPermission("customer:account:edit"), controller.AccountEdit)
	accountGroup.PUT("/:accountId/resetPwd", middlewares.HasPermission("customer:account:resetPwd"), controller.AccountResetPassword)
	accountGroup.DELETE("/:accountIds", middlewares.HasPermission("customer:account:remove"), controller.AccountRemove)
}

func InitPortalCustomerRouter(router *gin.RouterGroup) {
	group := router.Group("/portal/customer")
	group.POST("/login", controller.PortalCustomerLogin)
	group.Use(customermiddleware.CustomerAuthMiddleware())
	group.GET("/profile", controller.PortalCustomerProfile)
}
