package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("x-token")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"msg": "token不存在",
			})
			c.Abort()
			return
		}
		zap.S().Info("token msg", token)
		claims, err := ParseToken(token)
		zap.S().Info("claims:", claims)
		if err != nil {
			if err == TokenExpired {
				c.JSON(http.StatusOK, gin.H{
					"msg": "token授权已过期",
				})

			} else if err == TokenMalformed {
				c.JSON(http.StatusOK, gin.H{
					"msg": TokenMalformed.Error(),
				})

			} else if err == TokenNotValidYet {
				c.JSON(http.StatusOK, gin.H{
					"msg": TokenNotValidYet.Error(),
				})

			} else {
				c.JSON(http.StatusOK, gin.H{
					"msg": TokenInvalid.Error(),
				})

			}
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("userId", claims.ID)

		c.Next()
	}
}
