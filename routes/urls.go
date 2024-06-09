package routes

import (
	"GinChat/api"
	"GinChat/db"
	"GinChat/middleware"
	"GinChat/pkg/JWT"
	"GinChat/repository"
	"GinChat/service"
	"github.com/gin-gonic/gin"
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
	apiV1 := router.Group("api/v1")
	{
		//Authentication API
		auth := apiV1.Group("/auth", middleware.NotAuthorization())
		{
			auth.POST("", authAPI.Register)
			auth.PUT("", authAPI.Login)

		}
		user := apiV1.Group("/user", middleware.AuthorizationJWT(jwtAuth))
		{
			user.POST("")

		}

	}
	return router

}
