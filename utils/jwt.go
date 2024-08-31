package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateAccessToken(email, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"role":   role,
		"expire": time.Now().Add(accessExpire * time.Second).Unix()})

	return token.SignedString([]byte(SecretKey))
}

func GenerateRefreshToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"expire": time.Now().Add(refreshExpire * time.Second).Unix()})

	return token.SignedString([]byte(SecretKey))
}

func GenerateVerifyEmailToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"expire": time.Now().Add(verifyEmailExpire * time.Second).Unix()})
	return token.SignedString([]byte(SecretKey))
}

func GeneratePasswordToken(email, hashed_password string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":           email,
		"hashed_password": hashed_password,
		"expire":          time.Now().Add(passwordExpire * time.Second).Unix()})

	return token.SignedString([]byte(SecretKey))
}

func VerifyToken(token string) error {
	if token == "" {
		return errors.New("not authorized")
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(SecretKey), nil
	})

	if err != nil {
		return errors.New("could not parse token")
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		return errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return errors.New("invalid token claims")
	}

	expire := int64(claims["expire"].(float64))

	if time.Now().Unix() > expire {
		return errors.New("time has expired")
	}

	return nil
}
