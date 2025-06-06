package middleware

import (
	"bytes"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		blw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		c.Next()

		duration := time.Since(start)

		logger.Info("HTTP request",
			zap.String("method", c.Request.Method),
			zap.String("uri", c.Request.RequestURI),
			zap.Int("status", c.Writer.Status()),
			zap.Int("response_size", blw.body.Len()),
			zap.Duration("duration", duration),
		)
	}
}
