package agentRoutes

import (
	"baize/app/agent/controller"

	"github.com/gin-gonic/gin"
)

func InitFormActionRouter(router *gin.RouterGroup) {
	group := router.Group("/api/agent")
	group.POST("/form/submit", controller.SubmitAgentForm)
	group.POST("/action/execute", controller.ExecuteAgentAction)
}
