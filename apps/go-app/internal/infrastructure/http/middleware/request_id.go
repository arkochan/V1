package middleware

import "github.com/gin-gonic/gin"

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation to follow
		c.Next()
	}
}
