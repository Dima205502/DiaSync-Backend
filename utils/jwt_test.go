package utils

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

func TestGenerateAccessToken(t *testing.T) {
	var testCases = []struct {
		email  string
		role   string
		expire int64
	}{
		{"dmitrkozyrev2@gmail.com", "viewer", 0},
		{"mexasd123@gmail.com", "default", 0},
		{"romarkovet2004@gmail.com", "viewer", 0},
	}

	for _, tt := range testCases {

		tt.expire = time.Now().Add(accessExpire * time.Second).Unix()
		accessToken, err := GenerateAccessToken(tt.email, tt.role)

		if err != nil {
			t.Error(err)
		}

		parsedToken, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)

			if !ok {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(SecretKey), nil
		})

		if err != nil {
			t.Error("could not parse token")
		}

		tokenIsValid := parsedToken.Valid

		if !tokenIsValid {
			t.Error("invalid token")
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)

		if !ok {
			t.Error("invalid token claims")
		}

		_, ok = claims["expire"].(float64)

		if !ok {
			t.Fail()
		}

		token_expire := int64(claims["expire"].(float64))

		if token_expire-tt.expire > int64(5*time.Second) {
			t.Error("expire differense more then five second")
		}

		token_email, ok := claims["email"].(string)

		if !ok {
			t.Fail()
		}

		if token_email != tt.email {
			t.Errorf("got %s, want %s", token_email, tt.email)
		}

		token_role, ok := claims["role"].(string)

		if !ok {
			t.Fail()
		}

		if token_role != tt.role {
			t.Errorf("got %s, want %s", token_role, tt.role)
		}
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	expire := time.Now().Add(refreshExpire * time.Second).Unix()

	refreshToken, err := GenerateRefreshToken()

	if err != nil {
		t.Error(err.Error())
	}

	parsedToken, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(SecretKey), nil
	})

	if err != nil {
		t.Error("could not parse token")
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		t.Error("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		t.Error("invalid token claims")
	}

	_, ok = claims["expire"].(float64)

	if !ok {
		t.Fail()
	}

	token_expire := int64(claims["expire"].(float64))

	if token_expire-expire > int64(5*time.Second) {
		t.Error("expire differense more then five second")
	}

}

func TestGeneratePasswordToken(t *testing.T) {
	expire := time.Now().Add(passwordExpire * time.Second).Unix()

	passwordToken, err := GeneratePasswordToken("iopawndoiwqdno@yandex.ru", "ioadjioaun1i023hni12hj3nbi")

	if err != nil {
		t.Error(err.Error())
	}

	parsedToken, err := jwt.Parse(passwordToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(SecretKey), nil
	})

	if err != nil {
		t.Error("could not parse token")
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		t.Error("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		t.Error("invalid token claims")
	}

	_, ok = claims["expire"].(float64)

	if !ok {
		t.Fail()
	}

	token_expire := int64(claims["expire"].(float64))

	if token_expire-expire > int64(5*time.Second) {
		t.Error("expire differense more then five second")
	}

	email := claims["email"].(string)

	if email != "iopawndoiwqdno@yandex.ru" {
		t.Error("other email")
	}

	hashedPassword := claims["hashed_password"].(string)

	if hashedPassword != "ioadjioaun1i023hni12hj3nbi" {
		t.Error("other password")
	}
}

func TestGenerateVerifyEmailToken(t *testing.T) {
	expire := time.Now().Add(verifyEmailExpire * time.Second).Unix()

	verifyEmailToken, err := GenerateVerifyEmailToken("aopjdqonwd@gmail.com")

	if err != nil {
		t.Error(err.Error())
	}

	parsedToken, err := jwt.Parse(verifyEmailToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(SecretKey), nil
	})

	if err != nil {
		t.Error("could not parse token")
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		t.Error("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		t.Error("invalid token claims")
	}

	_, ok = claims["expire"].(float64)

	if !ok {
		t.Fail()
	}

	token_expire := int64(claims["expire"].(float64))

	if token_expire-expire > int64(5*time.Second) {
		t.Error("expire differense more then five second")
	}

	email := claims["email"].(string)

	if email != "aopjdqonwd@gmail.com" {
		t.Error("other email")
	}
}
