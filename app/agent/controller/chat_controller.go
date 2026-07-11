package controller

import (
	"baize/app/agent/request"
	"baize/app/agent/service"
	"baize/app/common/baize/baizeContext"
	customerService "baize/app/customer/service"
	freightService "baize/app/freight/service"
	"baize/app/utils/admin"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var chatService = service.GetChatService()

func CreateSession(c *gin.Context) {
	req := new(request.CreateSessionRequest)
	_ = c.ShouldBindJSON(req)
	c.JSON(200, chatService.CreateSession(currentUserID(c), req))
}

func ListModels(c *gin.Context) {
	c.JSON(200, chatService.ListModels())
}

func ListSessions(c *gin.Context) {
	c.JSON(200, chatService.ListSessions(currentUserID(c)))
}

func ListMessages(c *gin.Context) {
	sessionID, _ := strconv.ParseInt(c.Param("sessionId"), 10, 64)
	c.JSON(200, chatService.ListMessages(currentUserID(c), sessionID))
}

func UpdateSessionTitle(c *gin.Context) {
	sessionID, _ := strconv.ParseInt(c.Param("sessionId"), 10, 64)
	req := new(request.UpdateSessionTitleRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, gin.H{"message": "title is required"})
		return
	}
	if err := chatService.UpdateSessionTitle(currentUserID(c), sessionID, req); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "ok"})
}

func DeleteSession(c *gin.Context) {
	sessionID, _ := strconv.ParseInt(c.Param("sessionId"), 10, 64)
	if err := chatService.DeleteSession(currentUserID(c), sessionID); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "ok"})
}

func SendMessage(c *gin.Context) {
	req := new(request.SendMessageRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, gin.H{"message": "sessionId and message are required"})
		return
	}
	if strings.HasPrefix(c.Request.URL.Path, "/agent/chat/") {
		bzc := baizeContext.NewBaiZeContext(c)
		req.Source = "admin"
		req.OperatorID = bzc.GetCurrentUserId()
		req.OperatorName = bzc.GetCurrentUserName()
		req.CanManageAll = freightService.CanManageAllShipments(bzc.GetCurrentUser())
		req.Permissions = currentPermissions(bzc)
	} else if claims := customerClaimsFromHeader(c); claims != nil {
		req.Source = "customer"
		req.CustomerID = claims.CustomerId
		if customer := customerService.GetCustomerService().SelectCustomerById(claims.CustomerId); customer != nil { req.CustomerName = customer.CustomerName }
	}
	resp, err := chatService.Send(currentUserID(c), req)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func currentPermissions(bzc *baizeContext.BaiZeContext) []string {
	loginUser := bzc.GetCurrentLoginUser()
	if loginUser == nil || loginUser.User == nil {
		return nil
	}
	if admin.IsAdmin(loginUser.User.UserId) {
		return []string{"*:*:*"}
	}
	return loginUser.Permissions
}

func AnalyzeShipmentInSession(c *gin.Context) {
	sessionID, _ := strconv.ParseInt(c.Param("sessionId"), 10, 64)
	if file, err := c.FormFile("file"); err == nil {
		resp, err := chatService.AnalyzeShipmentFile(currentUserID(c), sessionID, file, c.PostForm("modelName"))
		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}
		c.JSON(200, resp)
		return
	}

	req := new(request.ShipmentAnalyzeRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, gin.H{"message": "invalid shipment analyze payload"})
		return
	}
	resp, err := chatService.AnalyzeShipment(currentUserID(c), sessionID, req)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func currentUserID(c *gin.Context) int64 {
	bzc := baizeContext.NewBaiZeContext(c)
	if bzc.GetCurrentLoginUser() != nil {
		return bzc.GetCurrentUserId()
	}
	if claims := customerClaimsFromHeader(c); claims != nil {
		if claims.AccountId > 0 {
			return -claims.AccountId
		}
		if claims.CustomerId > 0 {
			return -claims.CustomerId
		}
	}
	value, exists := c.Get("userId")
	if !exists {
		return 0
	}
	switch v := value.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	default:
		return 0
	}
}

func customerClaimsFromHeader(c *gin.Context) *customerService.CustomerClaims {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		return nil
	}
	token := strings.TrimSpace(auth)
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		token = strings.TrimSpace(token[7:])
	}
	if token == "" {
		return nil
	}
	claims, err := customerService.ParseCustomerToken(token)
	if err != nil {
		return nil
	}
	return claims
}
