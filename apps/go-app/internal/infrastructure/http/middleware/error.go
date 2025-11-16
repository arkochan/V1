package middleware

import "github.com/gin-gonic/gin"

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// Implementation to follow
	}
}
