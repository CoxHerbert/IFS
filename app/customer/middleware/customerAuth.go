package middleware

import (
	"baize/app/common/baize/baizeContext"
	"baize/app/constant/constants"
	"baize/app/customer/service"
	"strings"

	"github.com/gin-gonic/gin"
)

const CustomerClaimsKey = "customerClaims"

func CustomerAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bzc := baizeContext.NewBaiZeContext(c)
		authHeader := c.Request.Header.Get("Authorization")
		if len(authHeader) == 0 {
			bzc.InvalidToken()
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == constants.TokenPrefix) {
			bzc.InvalidToken()
			c.Abort()
			return
		}
		claims, err := service.ParseCustomerToken(parts[1])
		if err != nil {
			bzc.InvalidToken()
			c.Abort()
			return
		}
		c.Set(CustomerClaimsKey, claims)
		c.Next()
	}
}
