package server

import (
	"DiaSync/config"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

var clearPeriod time.Duration

func InitStorage(cfg config.Db) *Storage {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname)
	var err error

	DB, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic("coudn't connect to database")
	}

	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(5)

	CreateUsersTable(DB)
	CreateSessionsTable(DB)

	clearPeriod = cfg.ClearPeriod

	return &Storage{DB}
}

func CreateUsersTable(DB *sql.DB) {
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS Users(
	email TEXT PRIMARY KEY,
	password TEXT NOT NULL,
	role TEXT NOT NULL,
	verified BOOLEAN DEFAULT FALSE
	);`)

	if err != nil {
		panic(err.Error())
	}
}

func CreateSessionsTable(DB *sql.DB) {
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS Sessions(
	refresh_token TEXT PRIMARY KEY,
	user_email TEXT NOT NULL,
	deviceID TEXT NOT NULL,
	FOREIGN KEY (user_email) REFERENCES Users (email)
	);`)

	if err != nil {
		panic(err)
	}
}

func (s *Storage) Clear() {
	for {
		time.Sleep(clearPeriod)
		_, err := s.db.Exec(`DELETE FROM users WHERE verified=FALSE`)
		if err != nil {
			panic(err)
		}
	}
}
