package routes

import (
	"GinChat/api"
	"GinChat/db"
	_ "GinChat/docs"
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
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	//chat
	chatRepository repository.ChatRepository = repository.NewChatRepository(postDb, redisDb)
	chatService    service.ChatService       = service.NewChatService(chatRepository)
	chatAPI        api.ChatAPI               = api.NewChatAPI(chatService)

	//user
	userRepository repository.UserRepository = repository.NewUserRepository(postDb, redisDb)
	userService    service.UserService       = service.NewUserService(userRepository)
	userAPI        api.UserAPI               = api.NewUserAPI(userService)

	//jwt
	jwtAuth JWT.Jwt
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
			if err := v.RegisterValidation("name_validator", validators.NameValidate); err != nil {
				panic("validator issue in URLS")
			}
			if err := v.RegisterValidation("username_validator", validators.UsernameValidate); err != nil {
				panic("validator issue in URLS")
			}
			if err := v.RegisterValidation("image_validator", validators.ImageValidate); err != nil {
				panic("validator issue in URLS")
			}
		}

	}
	go websocketHandler.Manager.Start()
	//todo:need test
	// Serve Swagger documentation
	apiV1 := router.Group("api/v1")
	{
		//Authentication API
		auth := apiV1.Group("/auth", middleware.NotAuthorization())
		{
			auth.POST("", authAPI.Register)
			auth.PUT("", authAPI.Login, middleware.LoginAttemptCheck())
		}
		//user API
		user := apiV1.Group("/user", middleware.AuthorizationJWT(jwtAuth))
		{
			user.GET("/get-users/", userAPI.GetAllUsers)
			user.GET("/get-profile/:id/", userAPI.GetUserProfile)
			user.PUT("/profile-update/:id/", userAPI.ProfileUpdate)

		}
		//chat API
		chat := apiV1.Group("/chat", middleware.AuthorizationJWT(jwtAuth))
		{
			chat.GET("/ws/", chatAPI.ChatWs)
			chat.GET("get-rooms/", chatAPI.GetAllRooms)
			chat.POST("make-private/", chatAPI.MakePvChat)
			chat.POST("make-group/", chatAPI.MakeGroupChat)
			chat.POST("send-pv-message/", chatAPI.SendPvMessage)
			chat.POST("send-gp-message/", chatAPI.SendGpMessage)
		}

	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router

}
