package models

type LoginR struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
	DeviceID string `json:"device_id" binding:"required"`
	Role     string `json:"omitempty"`
}

type LogoutR struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ReplacementTokensR struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	DeviceID     string `json:"device_id" binding:"required"`
}

type ResetPasswordR struct {
	Email       string `binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type RepeatEmailVerifyR struct {
	Email string `binding:"required"`
}
