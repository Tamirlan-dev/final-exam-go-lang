package main

import (
	"github.com/Omar-Temirgali/final-exam-go-lang/config"
	"github.com/Omar-Temirgali/final-exam-go-lang/internal/controller"
	"github.com/Omar-Temirgali/final-exam-go-lang/internal/middleware"
	"github.com/Omar-Temirgali/final-exam-go-lang/internal/repository"
	"github.com/Omar-Temirgali/final-exam-go-lang/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db                *gorm.DB                     = config.SetupDatabaseConnection()
	userRepository    repository.UserRepository    = repository.NewUserRepository(db)
	articleRepository repository.ArticleRepository = repository.NewArticleRepository(db)
	jwtService        service.JWTService           = service.NewJWTService()
	userService       service.UserService          = service.NewUserService(userRepository)
	articleService    service.ArticleService       = service.NewArticleService(articleRepository)
	authService       service.AuthService          = service.NewAuthService(userRepository)
	authController    controller.AuthController    = controller.NewAuthController(authService, jwtService)
	articleController controller.ArticleController = controller.NewArticleController(articleService, jwtService)
	userController    controller.UserController    = controller.NewUserController(userService, jwtService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
	}

	articleRoutes := r.Group("api/articles", middleware.AuthorizeJWT(jwtService))
	{
		articleRoutes.GET("/", articleController.All)
		articleRoutes.GET("/:id", articleController.FindByID)
		articleRoutes.POST("/", articleController.Insert)
		articleRoutes.PUT("/:id", articleController.Update)
		articleRoutes.DELETE("/:id", articleController.Delete)
	}

	r.Run()
}
