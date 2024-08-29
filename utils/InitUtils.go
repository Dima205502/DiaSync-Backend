package utils

import "DiaSync/config"

var SecretKey string = "aaaaaaaa"
var access_expire_min int64 = 15
var refresh_expire_hour int64 = 24
var verify_email_expire_min int64 = 15

var app_password string = "aaaaaaaa"
var sender string = "aaaa@aaa.a"
var smtp_server string = "smtp.gmail.com"
var smtp_adr string = "smtp.gmail.com:587"

var server_adr string = "localhost:8080"

func Init(cfg config.Config) {
	server_adr = cfg.Server_adr
	InitEmail(cfg.Email)
	InitToken(cfg.Token)
}

func InitEmail(cfg config.Email) {
	app_password = cfg.App_password
	sender = cfg.Sender
	smtp_server = cfg.Smtp_server
	smtp_adr = cfg.Smtp_adr
}

func InitToken(cfg config.Token) {
	SecretKey = cfg.Secret_key
	access_expire_min = cfg.Access_expire_min
	refresh_expire_hour = cfg.Refresh_expire_hour
	verify_email_expire_min = cfg.Verify_email_expire_min
}
