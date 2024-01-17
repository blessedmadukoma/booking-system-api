package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// errorResponse handles the error response by using map[string]interface{} to return the error and it's message
func errorResponse(s string, err error) gin.H {
	if err != nil {
		return gin.H{"error": s + " -> " + err.Error()}
	}

	return gin.H{"error": s}
}

// rateLimitExceededResponse is a helper to send  a 429 To Many Requests
func (srv *Server) rateLimitExceededResponse(ctx *gin.Context) {
	ctx.JSON(http.StatusTooManyRequests, errorResponse("rate limit exceeded", nil))
}
