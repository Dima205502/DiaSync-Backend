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
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

type HttpServer struct {
	Server_adr  string        `json:"server_adr"`
	Timeout     time.Duration `json:"timeout"`
	IdleTimeout time.Duration `json:"idle_timeout"`
}

type Utils struct {
	Email `json:"email"`
	Token `json:"token"`
}

type Email struct {
	App_password string `json:"app_password"`
	Sender       string `json:"sender"`
	Smtp_server  string `json:"smtp_server"`
	Smtp_adr     string `json:"smtp_adr"`
}

type Token struct {
	Access_expire_min       int64  `json:"access_expire_min"`
	Refresh_expire_hour     int64  `json:"refresh_expire_hour"`
	Verify_email_expire_min int64  `json:"verify_email_expire_min"`
	Secret_key              string `json:"secret_key"`
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
