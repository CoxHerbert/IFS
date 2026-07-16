package controller

import (
	"baize/app/agent/request"
	"baize/app/agent/service"

	"github.com/gin-gonic/gin"
)

var agentConfigService = service.GetConfigService()

func GetOllamaConfig(c *gin.Context) {
	c.JSON(200, gin.H{"data": agentConfigService.GetOllamaConfig()})
}

func SaveOllamaConfig(c *gin.Context) {
	req := new(request.AgentOllamaConfigRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, gin.H{"message": "invalid agent config"})
		return
	}
	c.JSON(200, gin.H{"data": agentConfigService.SaveOllamaConfig(req), "message": "ok"})
}

func TestOllamaConfig(c *gin.Context) {
	req := new(request.AgentOllamaConfigRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, gin.H{"message": "invalid agent config"})
		return
	}
	if err := agentConfigService.TestOllamaConfig(req); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "ok"})
}
