package server

import (
	"DiaSync/config"
	"DiaSync/controller"
	"DiaSync/repository"
	"DiaSync/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(storage *Storage) *gin.Engine {
	authRepository := repository.NewAuthRepository(storage.db)
	authService := service.NewAuthService(authRepository)
	authController := controller.NewAuthController(authService)

	router := gin.New()

	auth := router.Group("/auth")

	{
		auth.POST("/signup", authController.Signup) // request --> email, password, role
		auth.POST("/verify-email", authController.VerifyEmail)
		auth.POST("/login", authController.Login)   // email password device_id
		auth.POST("/logout", authController.Logout) // refresh_token
		auth.POST("/replacement-token", authController.ReplacementTokens)
		auth.POST("/reset-password", authController.ResetPassword) // email, new_password
		auth.POST("/verify-newpassword", authController.VerifyNewPassword)
	}

	return router
}

func InitHttpServer(cfg config.Config, router http.Handler) *http.Server {
	server := &http.Server{
		Addr:         cfg.Server_adr,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	return server
}
