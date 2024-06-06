package middleware

import (
	"GinChat/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NoMethodHandler() gin.HandlerFunc {
	return func(request *gin.Context) {
		request.JSON(http.StatusMethodNotAllowed, utils.MethodNotAllowed)
	}
}

func NoRouteHandler() gin.HandlerFunc {
	return func(request *gin.Context) {
		request.JSON(http.StatusNotFound, utils.RouteNotDefined)
	}
}
