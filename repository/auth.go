package repository

import (
	"DiaSync/models"
	"DiaSync/utils"
	"database/sql"
	"errors"
)

type Authorization interface {
	CreateUser(models.User) error
	ValidateCredentials(string, string) (string, error)
	CreateSession(string, string, string) error
	GenerateTokens(string, string, string) (string, string, error)
	FindSession(string) (models.Session, error)
	DeleteRefreshToken(string) error
	FindUser(string) (models.User, error)
	VerifyEmail(string) error
	SetPassword(string, string) error
}

func NewAuthRepository(db *sql.DB) Authorization {
	return &AuthRepository{db}
}

type AuthRepository struct {
	db *sql.DB
}

func (s *AuthRepository) CreateUser(user models.User) error {
	hashedPassword := utils.HashPassword(user.Password)

	_, err := s.db.Exec("INSERT INTO Users (email, password, role) VALUES($1, $2, $3)", user.Email, hashedPassword, user.Role)

	return err
}

func (s *AuthRepository) ValidateCredentials(email, password string) (string, error) {
	row := s.db.QueryRow("SELECT password, role FROM Users WHERE email = $1;", email)

	var retrievedPassword, retrievedRole string
	err := row.Scan(&retrievedPassword, &retrievedRole)

	if err != nil {
		return "", err
	}

	passwordIsValid := utils.CheckPasswordHash(password, retrievedPassword)

	if !passwordIsValid {
		return "", errors.New("credential invalid")
	}

	return retrievedRole, nil
}

func (s *AuthRepository) CreateSession(refresh_token, user_email, deviceID string) error {
	_, err := s.db.Exec("INSERT INTO sessions (refresh_token, user_email, deviceID) VALUES($1, $2, $3)", refresh_token, user_email, deviceID)

	return err
}

func (s *AuthRepository) GenerateTokens(email, role, deviceID string) (string, string, error) {
	access_token, err := utils.GenerateAccessToken(email, role)

	if err != nil {
		return "", "", err
	}

	refresh_token, err := utils.GenerateRefreshToken()

	if err != nil {
		return "", "", err
	}

	err = s.CreateSession(refresh_token, email, deviceID)

	if err != nil {
		return "", "", err
	}

	return access_token, refresh_token, nil
}

func (s *AuthRepository) FindSession(refresh_token string) (models.Session, error) {
	row := s.db.QueryRow("SELECT * FROM Sessions WHERE refresh_token = $1;", refresh_token)

	var session models.Session

	err := row.Scan(&session.RefreshToken, &session.UserEmail, &session.DeviceID)

	return session, err
}

func (s *AuthRepository) DeleteRefreshToken(refresh_token string) error {
	_, err := s.db.Exec("DELETE FROM Sessions WHERE refresh_token = $1;", refresh_token)
	return err
}

func (s *AuthRepository) FindUser(email string) (models.User, error) {
	row := s.db.QueryRow("SELECT * FROM Users WHERE email = $1;", email)

	var user models.User
	var trash string
	err := row.Scan(&user.Email, &user.Password, &user.Role, &trash)

	return user, err
}

func (s *AuthRepository) VerifyEmail(email string) error {
	_, err := s.db.Exec("UPDATE Users SET verified=TRUE WHERE email=$1;", email)
	return err
}

func (s *AuthRepository) SetPassword(email, hashedPassword string) error {
	_, err := s.db.Exec("UPDATE Users SET password=$1 WHERE email=$2", hashedPassword, email)
	return err
}
