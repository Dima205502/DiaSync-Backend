package config

import (
	"encoding/json"
	"flag"
	"os"
	"time"
)

type Config struct {
	Utils      `json:"utils"`
	Db         `json:"db"`
	HttpServer `json:"httpServer"`
}

type Db struct {
	Host        string        `json:"host"`
	Port        int           `json:"port"`
	User        string        `json:"user"`
	Password    string        `json:"password"`
	Dbname      string        `json:"dbname"`
	ClearPeriod time.Duration `json:"clear_period"`
}

type HttpServer struct {
	ServerAdr   string        `json:"server_adr"`
	Timeout     time.Duration `json:"timeout"`
	IdleTimeout time.Duration `json:"idle_timeout"`
}

type Utils struct {
	Email `json:"email"`
	Token `json:"token"`
}

type Email struct {
	AppPassword string `json:"app_password"`
	Sender      string `json:"sender"`
	SmtpServer  string `json:"smtp_server"`
	SmtpAdr     string `json:"smtp_adr"`
}

type Token struct {
	AccessExpire      time.Duration `json:"access_expire"`
	RefreshExpire     time.Duration `json:"refresh_expire"`
	VerifyEmailExpire time.Duration `json:"verify_email_expire"`
	PasswordExpire    time.Duration `json:"password_expire"`
	SecretKey         string        `json:"secret_key"`
}

func Init() Config {
	path := flag.String("p", "", "path to config file")

	flag.Parse()

	file, err := os.ReadFile(*path)

	if err != nil {
		panic("Can't read config")
	}

	var cfg Config

	err = json.Unmarshal(file, &cfg)

	if err != nil {
		panic("Can't read config")
	}

	cfg.HttpServer.Timeout *= time.Second
	cfg.HttpServer.IdleTimeout *= time.Second

	return cfg
}
