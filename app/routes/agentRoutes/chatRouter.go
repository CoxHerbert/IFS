package agentRoutes

import (
	"baize/app/agent/controller"
	"baize/app/common/middlewares"

	"github.com/gin-gonic/gin"
)

func InitChatRouter(router *gin.RouterGroup) {
	group := router.Group("/agent-api/chat")
	group.GET("/models", controller.ListModels)
	group.POST("/session", controller.CreateSession)
	group.GET("/sessions", controller.ListSessions)
	group.GET("/session/:sessionId/messages", controller.ListMessages)
	group.PUT("/session/:sessionId/title", controller.UpdateSessionTitle)
	group.DELETE("/session/:sessionId", controller.DeleteSession)
	group.POST("/session/:sessionId/shipment-analyze", controller.AnalyzeShipmentInSession)
	group.POST("/send", controller.SendMessage)
}

func InitAgentManageChatRouter(router *gin.RouterGroup) {
	group := router.Group("/agent/chat")
	group.GET("/models", middlewares.HasPermission("ifs:agent:chat"), controller.ListModels)
	group.POST("/session", middlewares.HasPermission("ifs:agent:chat"), controller.CreateSession)
	group.GET("/sessions", middlewares.HasPermission("ifs:agent:chat"), controller.ListSessions)
	group.GET("/session/:sessionId/messages", middlewares.HasPermission("ifs:agent:chat"), controller.ListMessages)
	group.PUT("/session/:sessionId/title", middlewares.HasPermission("ifs:agent:chat"), controller.UpdateSessionTitle)
	group.DELETE("/session/:sessionId", middlewares.HasPermission("ifs:agent:chat"), controller.DeleteSession)
	group.POST("/session/:sessionId/shipment-analyze", middlewares.HasPermission("ifs:agent:chat"), controller.AnalyzeShipmentInSession)
	group.POST("/send", middlewares.HasPermission("ifs:agent:chat"), controller.SendMessage)
	group.POST("/form/submit", middlewares.HasPermission("ifs:agent:chat"), controller.SubmitAgentForm)

	configGroup := router.Group("/agent/config")
	configGroup.GET("/ollama", middlewares.HasPermission("ifs:agent:config"), controller.GetOllamaConfig)
	configGroup.PUT("/ollama", middlewares.HasPermission("ifs:agent:config"), controller.SaveOllamaConfig)
	configGroup.POST("/ollama/test", middlewares.HasPermission("ifs:agent:config"), controller.TestOllamaConfig)
}
