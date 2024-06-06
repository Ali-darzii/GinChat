package utils

import (
	"github.com/gin-gonic/gin"
)

var (
	MustNotAuthenticated    = gin.H{"status": false, "detail": "MUST_NOT_AUTHENTICATED", "error_code": "0"}
	AuthenticationRequired  = gin.H{"status": false, "detail": "Authentication_Required", "error_code": "1"}
	TokenIsExpiredOrInvalid = gin.H{"status": false, "detail": "Token_Expired_Or_Invalid", "error_code": "2"}
	MethodNotAllowed        = gin.H{"status": false, "detail": "Method_Not_Allowed", "error_code": "3"}
	RouteNotDefined         = gin.H{"status": false, "detail": "Route_Not_Defined", "error_code": "4"}
)
