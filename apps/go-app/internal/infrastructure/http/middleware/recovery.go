package middleware

import "github.com/gin-gonic/gin"

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// This can be replaced with gin.Recovery()
		c.Next()
	}
}
