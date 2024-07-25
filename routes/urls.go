package routes

import (
	"GinChat/api"
	"GinChat/db"
	"GinChat/middleware"
	"GinChat/pkg/JWT"
	"GinChat/repository"
	"GinChat/service"
	"GinChat/validators"
	"GinChat/websocketHandler"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	//db
	postDb  *gorm.DB      = db.ConnectPostgres()
	redisDb *redis.Client = db.ConnectRedis()

	//auth
	authRepository repository.AuthRepository = repository.NewAuthRepository(postDb, redisDb)
	authService    service.AuthService       = service.NewAuthService(authRepository)
	authAPI        api.AuthAPI               = api.NewAuthAPI(authService)

	//jwt
	jwtAuth JWT.Jwt

	//chat
	chatRepository repository.ChatRepository = repository.NewChatRepository(postDb, redisDb)
	chatService    service.ChatService       = service.NewChatService(chatRepository)
	chatAPI        api.ChatAPI               = api.NewChatAPI(chatService)
)

func Urls() *gin.Engine {
	router := gin.Default()
	//middlewares
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
	go websocketHandler.Manager.Start()

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
		//chat API
		chat := apiV1.Group("/chat", middleware.AuthorizationJWT(jwtAuth))
		{
			chat.GET("/ws/", chatAPI.ChatWs)
			chat.GET("get-users/", chatAPI.GetAllUsers)
			chat.GET("get-rooms/", chatAPI.GetAllRooms)
			chat.POST("make-chat/", chatAPI.MakePvChat)
		}

	}

	return router

}
