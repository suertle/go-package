package logrequest

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LogRequestHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capture the request body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("Failed to read request body: %v", err)
		}
		// Restore the request body to its original state
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Create a custom response writer to capture the response body
		customWriter := &CustomResponseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = customWriter

		logRequest := LogRequest{
			Method: c.Request.Method,
			Host:   c.Request.Host,
			Path:   c.Request.URL.Path,
			Query:  c.Request.URL.RawQuery,
			Header: headersToString(c.Request.Header),
			Body:   string(bodyBytes),
		}

		// handle panic, recover return status 500 and repanic
		defer func() {
			if r := recover(); r != nil {
				// Capture panic details in response
				panicMessage := fmt.Sprintf("%v", r)
				logRequest.Response = panicMessage
				logRequest.StatusCode = http.StatusInternalServerError

				if err := db.Create(&logRequest).Error; err != nil {
					log.Printf("Failed to log request: %v", err)
				}

				// hide error message in production
				if os.Getenv("APP_ENV") == "production" {
					panicMessage = "Internal Server Error"
				}

				// Ensure a proper response is sent after the panic
				c.JSON(http.StatusInternalServerError, gin.H{"error": panicMessage})

				// repanic
				panic(r)
			}
		}()

		// Process the request
		c.Next()

		// Get the response status code
		statusCode := c.Writer.Status()

		response := customWriter.body.String()
		// ignore 2xx response
		if statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices {
			response = ""
		}

		// Log the request details into the database
		logRequest.Response = response
		logRequest.StatusCode = statusCode

		if err := db.Create(&logRequest).Error; err != nil {
			log.Printf("Failed to log request: %v", err)
		}
	}
}

func headersToString(headers map[string][]string) string {
	var sb strings.Builder

	for key, values := range headers {
		sb.WriteString(fmt.Sprintf("%s: %s\n", key, strings.Join(values, ", ")))
	}

	return sb.String()
}
