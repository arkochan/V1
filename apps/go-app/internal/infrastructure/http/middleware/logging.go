package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func LoggingMiddleware(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		event := logger.Info()
		if statusCode >= 500 {
			event = logger.Error().Str("error", errorMessage)
		}

		event.
			Str("method", method).
			Str("path", path).
			Int("status_code", statusCode).
			Int("body_size", bodySize).
			Str("client_ip", clientIP).
			Dur("latency", latency).
			Msg("request")
	}
}
