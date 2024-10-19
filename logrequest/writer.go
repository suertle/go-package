package logrequest

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

// CustomResponseWriter wraps the gin.ResponseWriter to capture the response body
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)                  // Capture the response body
	return w.ResponseWriter.Write(b) // Write to the original response
}
