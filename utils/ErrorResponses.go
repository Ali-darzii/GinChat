package utils

import (
	"github.com/gin-gonic/gin"
)

var (
	MustNotAuthenticated    = gin.H{"status": false, "error_code": 0, "detail": "MUST_NOT_AUTHENTICATED"}
	AuthenticationRequired  = gin.H{"status": false, "error_code": 1, "detail": "Authentication_Required"}
	TokenIsExpiredOrInvalid = gin.H{"status": false, "error_code": 2, "detail": "Token_Expired_Or_Invalid"}
	MethodNotAllowed        = gin.H{"status": false, "error_code": 3, "detail": "Method_Not_Allowed"}
	RouteNotDefined         = gin.H{"status": false, "error_code": 4, "detail": "Route_Not_Defined"}
	BadFormat               = gin.H{"status": false, "error_code": 5, "detail": "Bad_Format"}
	UniqueField             = gin.H{"status": false, "error_code": 6, "detail": "Unique_Field_Required"}
	ObjectNotFound          = gin.H{"status": false, "error_code": 7, "detail": "Object_Not_Found"}
)
