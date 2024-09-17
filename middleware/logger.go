package middleware

import (
	"strconv"
	"time"

	"github.com/nitezs/sub2clash/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()
		latencyTime := endTime.Sub(startTime).Milliseconds()
		reqMethod := c.Request.Method
		reqURI := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		logger.Logger.Info(
			"Request",
			zap.Int("status", statusCode),
			zap.String("method", reqMethod),
			zap.String("uri", reqURI),
			zap.String("ip", clientIP),
			zap.String("latency", strconv.Itoa(int(latencyTime))+"ms"),
		)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Logger.Error(e)
			}
		}
	}
}
