package ledger

import (
	"baize/app/common/baize/baizeContext"
	"baize/app/setting"
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const claimsKey = "ledgerAccountId"

type tokenClaims struct {
	AccountID int64 `json:"accountId"`
	jwt.StandardClaims
}

func makeToken(accountID int64) (string, error) {
	// Keep ledger sessions active for one month, independent of the shorter
	// admin-console token lifetime configured in config.yaml.
	const ledgerTokenLifetime = 30 * 24 * time.Hour
	claims := tokenClaims{AccountID: accountID, StandardClaims: jwt.StandardClaims{Issuer: setting.Conf.TokenConfig.Issuer, ExpiresAt: time.Now().Add(ledgerTokenLifetime).Unix()}}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(setting.Conf.TokenConfig.Secret))
}

func parseToken(raw string) (*tokenClaims, error) {
	claims := new(tokenClaims)
	token, err := jwt.ParseWithClaims(raw, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(setting.Conf.TokenConfig.Secret), nil
	})
	if err != nil || !token.Valid || claims.AccountID <= 0 {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		parts := strings.SplitN(c.GetHeader("Authorization"), " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			baizeContext.NewBaiZeContext(c).InvalidToken()
			c.Abort()
			return
		}
		claims, err := parseToken(parts[1])
		if err != nil {
			baizeContext.NewBaiZeContext(c).InvalidToken()
			c.Abort()
			return
		}
		c.Set(claimsKey, claims.AccountID)
		c.Next()
	}
}

func accountID(c *gin.Context) (int64, bool) {
	id, ok := c.Get(claimsKey)
	if !ok {
		return 0, false
	}
	value, ok := id.(int64)
	return value, ok
}
