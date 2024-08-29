package controller

import (
	"DiaSync/models"
	"DiaSync/service"
	mock_service "DiaSync/service/mocks"
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestAuthController_Signup(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user models.User)

	var testCases = []struct {
		name                string
		inputBody           string
		inputUser           models.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email":"Dima", "password":"ddd", "role":"viewer"}`,
			inputUser: models.User{
				Email:    "Dima",
				Password: "ddd",
				Role:     "viewer",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user models.User) {
				s.EXPECT().CreateUser(user).Return(nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: ``,
		},
		{
			name:      "Incorrect Request",
			inputBody: `{"email":"Dima", "role":"viewer"}`,
			mockBehavior: func(s *mock_service.MockAuthorization, user models.User) {

			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"incorrect data"}`,
		},
		{
			name:      "Server error",
			inputBody: `{"email":"Dima", "password":"ddd", "role":"viewer"}`,
			inputUser: models.User{
				Email:    "Dima",
				Password: "ddd",
				Role:     "viewer",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user models.User) {
				s.EXPECT().CreateUser(user).Return(errors.New("Server error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"couldn't create the user"}`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			tt.mockBehavior(auth, tt.inputUser)

			authService := service.Authorization(auth)
			authController := NewAuthController(authService)

			r := gin.New()
			r.POST("/signup", authController.Signup)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/signup", bytes.NewBufferString(tt.inputBody))

			r.ServeHTTP(w, req)

			if tt.expectedStatusCode != w.Code {
				t.Errorf("got = %d expected = %d", w.Code, tt.expectedStatusCode)
			}

			if tt.expectedRequestBody != w.Body.String() {
				t.Errorf("got = %s expected = %s", w.Body.String(), tt.expectedRequestBody)
			}
		})
	}
}

func TestAuthController_Login(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user models.LoginR)

	var testCases = []struct {
		name                string
		inputBody           string
		inputUser           models.LoginR
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email":"Dima", "password":"ddd", "device_id":"DDD"}`,
			inputUser: models.LoginR{
				Email:    "Dima",
				Password: "ddd",
				DeviceID: "DDD",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user models.LoginR) {
				s.EXPECT().GenerateTokens(user).Return("asdasdads", "sadasfasfda", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"access_token":"asdasdads","refresh_token":"sadasfasfda"}`,
		},
		{
			name:      "Incorrect request",
			inputBody: `{"email":"Dima", "password":"ddd"}`,
			inputUser: models.LoginR{
				Email:    "Dima",
				Password: "ddd",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user models.LoginR) {

			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"incorrect data"}`,
		},
		{
			name:      "Server error",
			inputBody: `{"email":"Dima", "password":"ddd", "device_id":"DDD"}`,
			inputUser: models.LoginR{
				Email:    "Dima",
				Password: "ddd",
				DeviceID: "DDD",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user models.LoginR) {
				s.EXPECT().GenerateTokens(user).Return("", "", errors.New("Server error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"couldn't generate tokens"}`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			tt.mockBehavior(auth, tt.inputUser)

			authService := service.Authorization(auth)
			authController := NewAuthController(authService)

			r := gin.New()
			r.POST("/login", authController.Login)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(tt.inputBody))

			r.ServeHTTP(w, req)

			if tt.expectedStatusCode != w.Code {
				t.Errorf("got = %d expected = %d", w.Code, tt.expectedStatusCode)
			}

			if tt.expectedRequestBody != w.Body.String() {
				t.Errorf("got = %s expected = %s", w.Body.String(), tt.expectedRequestBody)
			}
		})
	}
}

func TestAuthController_Logout(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user models.LogoutR)

	var testCases = []struct {
		name                string
		inputBody           string
		inputUser           models.LogoutR
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"refresh_token":"asdasdasfmkm"}`,
			inputUser: models.LogoutR{
				RefreshToken: "asdasdasfmkm",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, request models.LogoutR) {
				s.EXPECT().DeleteSession(request).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: ``,
		},
		{
			name: "Incorrect request",
			mockBehavior: func(s *mock_service.MockAuthorization, request models.LogoutR) {

			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"incorrect data"}`,
		},
		{
			name:      "Server error",
			inputBody: `{"refresh_token":"asdasdasfmkm"}`,
			inputUser: models.LogoutR{
				RefreshToken: "asdasdasfmkm",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, request models.LogoutR) {
				s.EXPECT().DeleteSession(request).Return(errors.New("Server error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"couldn't delete session"}`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			tt.mockBehavior(auth, tt.inputUser)

			authService := service.Authorization(auth)
			authController := NewAuthController(authService)

			r := gin.New()
			r.POST("/logout", authController.Logout)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/logout", bytes.NewBufferString(tt.inputBody))

			r.ServeHTTP(w, req)

			if tt.expectedStatusCode != w.Code {
				t.Errorf("got = %d expected = %d", w.Code, tt.expectedStatusCode)
			}

			if tt.expectedRequestBody != w.Body.String() {
				t.Errorf("got = %s expected = %s", w.Body.String(), tt.expectedRequestBody)
			}
		})
	}
}

func TestAuthController_ReplacementToken(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user models.ReplacementTokensR)

	var testCases = []struct {
		name                string
		inputBody           string
		inputUser           models.ReplacementTokensR
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"refresh_token":"asdasdasfmkm", "device_id":"DDD"}`,
			inputUser: models.ReplacementTokensR{
				RefreshToken: "asdasdasfmkm",
				DeviceID:     "DDD",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, request models.ReplacementTokensR) {
				s.EXPECT().ReplacementTokens(request).Return("sfdfadfdsaf", "ojoiewjeq", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"access_token":"sfdfadfdsaf","refresh_token":"ojoiewjeq"}`,
		},
		{
			name:      "Incorrect request",
			inputBody: `{"refresh_token":"asdasdasfmkm"}`,
			inputUser: models.ReplacementTokensR{
				RefreshToken: "asdasdasfmkm",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, request models.ReplacementTokensR) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"incorrect data"}`,
		},
		{
			name:      "Server error",
			inputBody: `{"refresh_token":"asdasdasfmkm", "device_id":"DDD"}`,
			inputUser: models.ReplacementTokensR{
				RefreshToken: "asdasdasfmkm",
				DeviceID:     "DDD",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, request models.ReplacementTokensR) {
				s.EXPECT().ReplacementTokens(request).Return("", "", errors.New("Server error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"couldn't replacement tokens"}`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			tt.mockBehavior(auth, tt.inputUser)

			authService := service.Authorization(auth)
			authController := NewAuthController(authService)

			r := gin.New()
			r.POST("/replacement-token", authController.ReplacementTokens)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/replacement-token", bytes.NewBufferString(tt.inputBody))

			r.ServeHTTP(w, req)

			if tt.expectedStatusCode != w.Code {
				t.Errorf("got = %d expected = %d", w.Code, tt.expectedStatusCode)
			}

			if tt.expectedRequestBody != w.Body.String() {
				t.Errorf("got = %s expected = %s", w.Body.String(), tt.expectedRequestBody)
			}
		})
	}
}
