package routes

import (
	"GinChat/api"
	"GinChat/db"
	"GinChat/middleware"
	"GinChat/pkg/JWT"
	"GinChat/repository"
	"GinChat/service"
	"GinChat/validators"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var (
	postDb         *gorm.DB                  = db.ConnectPostgres()
	authRepository repository.AuthRepository = repository.NewAuthRepository(postDb)
	authService    service.AuthService       = service.NewAuthService(authRepository)
	authAPI        api.AuthAPI               = api.NewAuthAPI(authService)
	jwtAuth        JWT.Jwt
)

func Urls() *gin.Engine {
	router := gin.Default()
	//middleware
	{
		router.Use(middleware.CORSMiddleware())
		router.NoRoute(middleware.NoRouteHandler())
		router.HandleMethodNotAllowed = true
		router.NoMethod(middleware.NoMethodHandler())
	}
	//register validations
	{
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			if err := v.RegisterValidation("phone_validator", validators.PhoneNoValidate); err != nil {
				panic("validator issue in URLS")
			}
			if err := v.RegisterValidation("name_validator", validators.NameValidator); err != nil {
				panic("validator issue in URLS")
			}
			if err := v.RegisterValidation("username_validator", validators.UsernameValidate); err != nil {
				panic("validator issue in URLS")
			}
		}

	}

	apiV1 := router.Group("api/v1")
	{
		//Authentication API
		auth := apiV1.Group("/auth", middleware.NotAuthorization())
		{
			auth.POST("", authAPI.Register)
			auth.PUT("", authAPI.Login, middleware.LoginAttemptCheck())

		}
		user := apiV1.Group("/user", middleware.AuthorizationJWT(jwtAuth))
		{
			user.POST("")

		}

	}
	return router

}
