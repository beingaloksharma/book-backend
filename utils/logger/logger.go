package logger

import (
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	logrus.SetOutput(gin.DefaultWriter)
	logrus.SetLevel(logrus.InfoLevel)
	// ReportCaller adds the file and line number to the log
	logrus.SetReportCaller(true)
}

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

// LogError logs the error with stack trace context and sends a JSON response
func LogError(c *gin.Context, code int, err error, msg string) {
	if err != nil {
		// Capture the caller (1 stack frame up)
		pc, file, line, ok := runtime.Caller(1)
		fields := logrus.Fields{
			"status_code": code,
			"path":        c.Request.URL.Path,
			"error":       err.Error(),
		}
		if ok {
			fields["file"] = file
			fields["line"] = line
			fields["func"] = runtime.FuncForPC(pc).Name()
		}
		logrus.WithFields(fields).Error("API Error Encountered")
		c.Error(err) // Attach to Gin context for middleware logging if needed
	}

	response := gin.H{"error": msg}
	if msg == "" && err != nil {
		response["error"] = err.Error()
	}

	c.JSON(code, response)
}
