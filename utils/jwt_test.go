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

		tt.expire = time.Now().Add(time.Minute * time.Duration(access_expire_min)).Unix()
		accessToken, err := GenerateAccessToken(tt.email, tt.role)

		if err != nil {
			t.Error(err)
		}

		if accessToken == "" {
			t.Error("token is empty")
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
	expire := time.Now().Add(time.Minute * time.Duration(access_expire_min)).Unix()

	refreshToken, err := GenerateRefreshToken()

	if err != nil {
		t.Error(err.Error())
	}

	if refreshToken == "" {
		t.Error("token is empty")
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

func TestVerifyToken(t *testing.T) {
	var testCases = []string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmUiOjE3MjQ0MTM2NDZ9.OJqTLEfG6j8F1f7lMONjegc3hv1nNiX1yhMT-w1KUTc",
		"skdfskjfdh", "", "@"}

	t.Log("all tokens in tests are invalid")

	for _, tt := range testCases {
		err := VerifyToken(tt)

		if err == nil {
			t.Error("token is valid!")
		}
	}
}
