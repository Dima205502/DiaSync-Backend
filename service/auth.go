package service

import (
	"DiaSync/models"
	"DiaSync/repository"
	"DiaSync/utils"
	"errors"

	"github.com/golang-jwt/jwt"
)

//go:generate mockgen -source=auth.go -destination=mocks/mock.go
type Authorization interface {
	CreateUser(models.User) error
	GenerateTokens(models.LoginR) (string, string, error)
	DeleteSession(models.LogoutR) error
	ReplacementTokens(models.ReplacementTokensR) (string, string, error)
	VerifyEmail(string) error
	ResetPassword(models.ResetPasswordR) error
	VerifyNewPassword(string) error
	RepeatEmailVerify(string) error
}

func NewAuthService(authRepository repository.Authorization) Authorization {
	return &AuthService{authRepository}
}

type AuthService struct {
	AuthRepository repository.Authorization
}

func (as *AuthService) CreateUser(user models.User) error {
	err := as.AuthRepository.CreateUser(user)

	if err != nil {
		return err
	}

	verifyEmailToken, err := utils.GenerateVerifyEmailToken(user.Email)

	if err != nil {
		return err
	}

	return utils.SendVerifyTokenMail(user.Email, verifyEmailToken)
}

func (as *AuthService) GenerateTokens(userInfo models.LoginR) (string, string, error) {
	var err error
	userInfo.Role, err = as.AuthRepository.ValidateCredentials(userInfo.Email, userInfo.Password)

	if err != nil {
		return "", "", nil
	}

	return as.AuthRepository.GenerateTokens(userInfo.Email, userInfo.Role, userInfo.DeviceID)
}

func (as *AuthService) DeleteSession(request models.LogoutR) error {
	_, err := as.AuthRepository.FindSession(request.RefreshToken)

	if err != nil {
		return err
	}

	return as.AuthRepository.DeleteRefreshToken(request.RefreshToken)
}

func (as *AuthService) ReplacementTokens(request models.ReplacementTokensR) (string, string, error) {
	session, err := as.AuthRepository.FindSession(request.RefreshToken)

	if err != nil {
		return "", "", err
	}

	err = as.AuthRepository.DeleteRefreshToken(request.RefreshToken)

	if err != nil {
		return "", "", err
	}

	if session.DeviceID != request.DeviceID {
		return "", "", errors.New("invalid data")
	}

	user, err := as.AuthRepository.FindUser(session.UserEmail)

	if err != nil {
		return "", "", err
	}

	return as.AuthRepository.GenerateTokens(user.Email, user.Role, request.DeviceID)
}

func (as *AuthService) VerifyEmail(token string) error {
	err := utils.VerifyToken(token)

	if err != nil {
		return err
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(utils.SecretKey), nil
	})

	if err != nil {
		return err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return errors.New("invalid token claims")
	}

	email := claims["email"].(string)

	return as.AuthRepository.VerifyEmail(email)
}

func (as *AuthService) ResetPassword(request models.ResetPasswordR) error {
	hashedNewPassword := utils.HashPassword(request.NewPassword)

	newPasswordToken, err := utils.GeneratePasswordToken(request.Email, hashedNewPassword)

	if err != nil {
		return err
	}

	return utils.SendNewPasswordEmail(request.Email, newPasswordToken)
}

func (as *AuthService) VerifyNewPassword(token string) error {
	err := utils.VerifyToken(token)

	if err != nil {
		return err
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(utils.SecretKey), nil
	})

	if err != nil {
		return err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return errors.New("invalid token claims")
	}

	email := claims["email"].(string)
	hashedNewPassword := claims["hashed_password"].(string)

	return as.AuthRepository.SetPassword(email, hashedNewPassword)
}

func (as *AuthService) RepeatEmailVerify(email string) error {
	verifyEmailToken, err := utils.GenerateVerifyEmailToken(email)

	if err != nil {
		return err
	}

	return utils.SendVerifyTokenMail(email, verifyEmailToken)
}
