package models

type User struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
	Role     string `binding:"required"`
}

type Session struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	UserEmail    string `json:"user_email" binding:"required"`
	DeviceID     string `json:"device_id" binding:"required"`
}
