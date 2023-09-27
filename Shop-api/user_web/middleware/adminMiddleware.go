package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// AdminAuth 权限认证
// 1代表admin
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		value, _ := c.Get("claims")
		zap.S().Info("claims", value)
		claims := value.(*UserBasicClaims)
		if claims.AuthorityId != 2 {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "权限不足",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
