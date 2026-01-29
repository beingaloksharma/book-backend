package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()
		// Process request
		c.Next()
		// Stop timer
		endTime := time.Now()
		latency := endTime.Sub(startTime)
		// Get Status
		statusCode := c.Writer.Status()
		// Get Client IP
		clientIP := c.ClientIP()
		// Get Method and Path
		method := c.Request.Method
		path := c.Request.URL.Path
		// Create a logger entry with fields
		entry := logrus.WithFields(logrus.Fields{
			"status_code": statusCode,
			"latency":     latency, // Logrus formats durations automatically
			"client_ip":   clientIP,
			"method":      method,
			"path":        path,
			"user_agent":  c.Request.UserAgent(),
		})
		// Log based on status code
		if len(c.Errors) > 0 {
			// If there are internal Gin errors, log as Error
			entry.Error(c.Errors.String())
		} else if statusCode >= 500 {
			entry.Error("Internal Server Error")
		} else if statusCode >= 400 {
			entry.Warn("Bad Request / Client Error")
		} else {
			entry.Info("Request Processed")
		}
	}
}
