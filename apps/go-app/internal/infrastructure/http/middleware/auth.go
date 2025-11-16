package middleware

import "github.com/gin-gonic/gin"

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation to follow
		c.Next()
	}
}
