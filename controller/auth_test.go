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

func TestAuthController_ResetPassword(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user models.ResetPasswordR)

	var testCases = []struct {
		name                string
		inputBody           string
		inputUser           models.ResetPasswordR
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email":"asdasdasfmkm@gmail.com", "new_password":"III"}`,
			inputUser: models.ResetPasswordR{
				Email:       "asdasdasfmkm@gmail.com",
				NewPassword: "III",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, request models.ResetPasswordR) {
				s.EXPECT().ResetPassword(request).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "",
		},
		{
			name:      "Bad request",
			inputBody: `{"email":"asdaskap;oqpm@yandex.ru"}`,
			inputUser: models.ResetPasswordR{
				Email: "asdaskap;oqpm@yandex.ru",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, request models.ResetPasswordR) {

			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"incorrect data"}`,
		},
		{
			name:      "Internal server error",
			inputBody: `{"email":"soojwqdmqwpjmlme@yandex.ru", "new_password":"III"}`,
			inputUser: models.ResetPasswordR{
				Email:       "soojwqdmqwpjmlme@yandex.ru",
				NewPassword: "III",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, request models.ResetPasswordR) {
				s.EXPECT().ResetPassword(request).Return(errors.New("kadkolokad"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: ``,
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
			r.POST("/reset-password", authController.ResetPassword)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/reset-password", bytes.NewBufferString(tt.inputBody))

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

func TestAuthController_RepeatEmailVerify(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, email string)

	var testCases = []struct {
		name                string
		inputBody           string
		inputUser           models.RepeatEmailVerifyR
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email":"asdasdasfmkm@gmail.com"}`,
			inputUser: models.RepeatEmailVerifyR{
				Email: "asdasdasfmkm@gmail.com",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, email string) {
				s.EXPECT().RepeatEmailVerify(email).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "",
		},
		{
			name:      "Bad request",
			inputBody: `{"emails":"asdasdasfmkm@gmail.com"}`,
			mockBehavior: func(s *mock_service.MockAuthorization, email string) {

			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"incorrect data"}`,
		},
		{
			name:      "Internal server error",
			inputBody: `{"email":"asdasdasfmkm@gmail.com"}`,
			inputUser: models.RepeatEmailVerifyR{
				Email: "asdasdasfmkm@gmail.com",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, email string) {
				s.EXPECT().RepeatEmailVerify(email).Return(errors.New("couldn't create token"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"couldn't create token"}`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			tt.mockBehavior(auth, tt.inputUser.Email)

			authService := service.Authorization(auth)
			authController := NewAuthController(authService)

			r := gin.New()
			r.POST("/repeat-verify-email", authController.RepeatEmailVerify)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/repeat-verify-email", bytes.NewBufferString(tt.inputBody))

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

func TestAuthController_VerifyEmail(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, verifyEmailToken string)

	var testCases = []struct {
		name                string
		verifyEmailToken    string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:             "OK",
			verifyEmailToken: "JWONQW132NJ12NO213.O123NOJN1K.KONJKIN3O231NOL",
			mockBehavior: func(s *mock_service.MockAuthorization, verifyEmailToken string) {
				s.EXPECT().VerifyEmail(verifyEmailToken).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: ``,
		},
		{
			name:             "Bad request",
			verifyEmailToken: "OPKI13O12KK1N3M1L.ED23NKJ1K.KOO12UI54JHKB",
			mockBehavior: func(s *mock_service.MockAuthorization, verifyEmailToken string) {
				s.EXPECT().VerifyEmail(verifyEmailToken).Return(errors.New("invalid token"))
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid token"}`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			tt.mockBehavior(auth, tt.verifyEmailToken)

			authService := service.Authorization(auth)
			authController := NewAuthController(authService)

			r := gin.New()
			r.POST("/verify-email", authController.VerifyEmail)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/verify-email?token="+tt.verifyEmailToken, nil)

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

func TestAuthController_VerifyNewPassword(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, newPasswordToken string)

	var testCases = []struct {
		name                string
		newPasswordToken    string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:             "OK",
			newPasswordToken: "JWONQW132NJ12NO213.O123NOJN1K.KONJKIN3O231NOL",
			mockBehavior: func(s *mock_service.MockAuthorization, newPasswordToken string) {
				s.EXPECT().VerifyNewPassword(newPasswordToken).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: ``,
		},
		{
			name:             "Bad request",
			newPasswordToken: "JWONQW132NJ12NO213.O123NOJN1K.KONJKIN3O231NOL",
			mockBehavior: func(s *mock_service.MockAuthorization, newPasswordToken string) {
				s.EXPECT().VerifyNewPassword(newPasswordToken).Return(errors.New("invalid token"))
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid token"}`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			tt.mockBehavior(auth, tt.newPasswordToken)

			authService := service.Authorization(auth)
			authController := NewAuthController(authService)

			r := gin.New()
			r.POST("/verify-newpassword", authController.VerifyNewPassword)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/verify-newpassword?token="+tt.newPasswordToken, nil)

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
