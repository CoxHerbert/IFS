package systemRouter

import (
	"baize/app/system/controller/loginController"
	"github.com/gin-gonic/gin"
)

func InitGetUser(router *gin.RouterGroup) {
	router.GET("/getInfo", loginController.GetInfo)
	router.GET("/getRouters", loginController.GetRouters)
}

func InitLoginRouter(router *gin.RouterGroup) {
	router.GET("/captchaImage", loginController.GetCode)
	router.GET("/portal/captchaImage", loginController.GetCode)
	router.POST("/login", loginController.Login)
	router.POST("/logout", loginController.Logout)
}
