package handler

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"log/slog"
	"strconv"
	"time"
)

func Logging(c *gin.Context) {
	start := time.Now()
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
	}

	// Restore the request body for the next handlers
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	c.Next()
	l := slog.With(
		"statusCode", c.Writer.Status(),
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
		"params", c.Request.URL.Query(),
		"body", string(bodyBytes),
		"duration", time.Since(start).String())
	if c.Writer.Status() >= 500 {
		l.Error("response")
	}
	l.Info("response status")
}

func GetIntQueryParam(c *gin.Context, param string, defaultValue int) (int, error) {
	valueStr := c.Query(param)
	if valueStr == "" {
		return defaultValue, nil
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid value for %s: %v", param, err)
	}
	return value, nil
}

func ParseDateParam(c *gin.Context, param, defaultValue string) (time.Time, error) {
	dateStr := c.Query(param)
	if dateStr == "" {
		return time.Parse(time.DateOnly, defaultValue)
	}
	parsedDate, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date for %s: %v", param, err)
	}
	return parsedDate, nil
}
