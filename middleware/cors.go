package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(request *gin.Context) {
		request.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		request.Writer.Header().Set("Access-Control-Allow-Credentials", "*")
		request.Writer.Header().Set("Access-Control-Allow-Headers", "true")
		request.Writer.Header().Set("Access-Control-Allow-", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Catch-Control, X-Request-With")
		request.Writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS,POST, PUT, DELETE")
		if request.Request.Method == "OPTIONS" {
			request.AbortWithStatus(http.StatusNoContent)
		}
		request.Next()
	}
}
