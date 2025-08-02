package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func errorMessage(c *gin.Context, status int, message string, headers http.Header) {
	jsonWithHeaders(c, status, gin.H{"ErrorTrace": message}, headers)
}

func notFound(c *gin.Context) {
	message := "The requested resource could not be found"
	errorMessage(c, http.StatusNotFound, message, nil)
}

func methodNotAllowed(c *gin.Context) {
	message := fmt.Sprintf("The %s method is not supported for this resource", c.Request.Method)
	errorMessage(c, http.StatusMethodNotAllowed, message, nil)
}

func jsonWithHeaders(c *gin.Context, status int, data any, headers http.Header) {
	for key, value := range headers {
		c.Writer.Header()[key] = value
	}
	c.JSON(status, data)
}
