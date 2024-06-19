package middleware

import (
	"GinChat/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func LoginAttemptCheck() gin.HandlerFunc {
	return func(request *gin.Context) {
		loginAttempt, err := request.Cookie("loginAttempt")
		if err != nil {
			request.SetCookie("loginAttempt",
				"0",
				5400,
				"/",
				"localhost",
				false,
				false,
			)

		} else {
			attemptCheck, err := strconv.Atoi(loginAttempt)
			if err != nil {
				request.AbortWithStatusJSON(http.StatusBadRequest, utils.SomethingWentWrong)
				return
			}
			if attemptCheck >= 12 {
				request.AbortWithStatusJSON(http.StatusTooManyRequests, utils.TooManyLoginRequest)
				return
			} else {
				attemptCheck += 1
				loginAttempt = strconv.Itoa(attemptCheck)
				request.SetCookie("loginAttempt",
					loginAttempt,
					5400,
					"/",
					"localhost",
					false,
					false,
				)

			}
		}
		return
	}
}
