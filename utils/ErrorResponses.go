package utils

type ErrorResponse struct {
	Status    bool   `json:"status"`
	ErrorCode uint   `json:"error_code"`
	Detail    string `json:"detail"`
}

var (
	MustNotAuthenticated = ErrorResponse{
		Status:    false,
		ErrorCode: 0,
		Detail:    "MUST_NOT_AUTHENTICATED",
	}
	AuthenticationRequired = ErrorResponse{
		Status:    false,
		ErrorCode: 1,
		Detail:    "Authentication_Required",
	}
	TokenIsExpiredOrInvalid = ErrorResponse{
		Status:    false,
		ErrorCode: 2,
		Detail:    "Token_Expired_Or_Invalid",
	}
	MethodNotAllowed = ErrorResponse{
		Status:    false,
		ErrorCode: 3,
		Detail:    "Method_Not_Allowed",
	}
	RouteNotDefined = ErrorResponse{
		Status:    false,
		ErrorCode: 4,
		Detail:    "Route_Not_Defined",
	}
	BadFormat = ErrorResponse{
		Status:    false,
		ErrorCode: 5,
		Detail:    "Bad_Format",
	}
	ObjectNotFound = ErrorResponse{
		Status:    false,
		ErrorCode: 6,
		Detail:    "Object_Not_Found",
	}
	ExpiredTimeBlocked = ErrorResponse{
		Status:    false,
		ErrorCode: 7,
		Detail:    "Too_Many_Token_Request",
	}
	SomethingWentWrong = ErrorResponse{
		Status:    false,
		ErrorCode: 8,
		Detail:    "We_Don't_Know_What_Happened",
	}
	TooManyLoginRequest = ErrorResponse{
		Status:    false,
		ErrorCode: 9,
		Detail:    "Ip_Banned_For_15_min",
	}
	CanNotReachWsConnection = ErrorResponse{
		Status:    false,
		ErrorCode: 10,
		Detail:    "Can_Not_Reach_Ws_Connection",
	}
	UserNameIsTaken = ErrorResponse{
		Status:    false,
		ErrorCode: 11,
		Detail:    "User_Name_Is_Taken",
	}
	NameFieldRequired = ErrorResponse{
		Status:    false,
		ErrorCode: 12,
		Detail:    "Name_Field_Required_For_Register",
	}
)
