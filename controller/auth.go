package controller

import (
	"DiaSync/models"
	"DiaSync/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Authorization interface {
	Signup(*gin.Context)
	Login(*gin.Context)
	Logout(*gin.Context)
	ReplacementTokens(*gin.Context)
	VerifyEmail(*gin.Context)
	ResetPassword(*gin.Context)
	VerifyNewPassword(*gin.Context)
}

func NewAuthController(authService service.Authorization) Authorization {
	return &AuthController{authService}
}

type AuthController struct {
	authService service.Authorization
}

func (ac *AuthController) Signup(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "incorrect data"})
		return
	}

	err = ac.authService.CreateUser(user)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "couldn't create the user"})
		return
	}

	context.Status(http.StatusCreated)
}

func (ac *AuthController) Login(context *gin.Context) {
	var userInfo models.LoginR

	err := context.ShouldBindJSON(&userInfo)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "incorrect data"})
		return
	}

	access_token, refresh_token, err := ac.authService.GenerateTokens(userInfo)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "couldn't generate tokens"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"access_token": access_token,
		"refresh_token": refresh_token})
}

func (ac *AuthController) Logout(context *gin.Context) {
	var request models.LogoutR

	err := context.ShouldBindJSON(&request)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "incorrect data"})
		return
	}

	err = ac.authService.DeleteSession(request)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "couldn't delete session"})
		return
	}

	context.Status(http.StatusOK)
}

func (ac *AuthController) ReplacementTokens(context *gin.Context) {
	var request models.ReplacementTokensR

	err := context.ShouldBindJSON(&request)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "incorrect data"})
		return
	}

	access_token, refresh_token, err := ac.authService.ReplacementTokens(request)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"access_token": access_token, "refresh_token": refresh_token})
}

func (ac *AuthController) VerifyEmail(context *gin.Context) {
	verifyToken := context.Query("token")

	err := ac.authService.VerifyEmail(verifyToken)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.Status(http.StatusOK)
}

func (ac *AuthController) ResetPassword(context *gin.Context) {
	var request models.ResetPasswordR

	err := context.ShouldBindJSON(&request)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "incorrect data"})
		return
	}

	err = ac.authService.ResetPassword(request)

	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}

	context.Status(http.StatusOK)
}

func (ac *AuthController) VerifyNewPassword(context *gin.Context) {
	token := context.Query("token")

	err := ac.authService.VerifyNewPassword(token)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.Status(http.StatusOK)
}
