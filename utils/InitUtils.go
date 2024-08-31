package utils

import (
	"DiaSync/config"
	"time"
)

var SecretKey string
var accessExpire time.Duration
var refreshExpire time.Duration
var verifyEmailExpire time.Duration
var passwordExpire time.Duration

var appPassword string
var sender string
var smtpServer string
var smtpAdr string

var serverAdr string

func Init(cfg config.Config) {
	serverAdr = cfg.ServerAdr
	InitEmail(cfg.Email)
	InitToken(cfg.Token)
}

func InitEmail(cfg config.Email) {
	appPassword = cfg.AppPassword
	sender = cfg.Sender
	smtpServer = cfg.SmtpServer
	smtpAdr = cfg.SmtpAdr
}

func InitToken(cfg config.Token) {
	SecretKey = cfg.SecretKey
	accessExpire = cfg.AccessExpire
	refreshExpire = cfg.RefreshExpire
	verifyEmailExpire = cfg.VerifyEmailExpire
	passwordExpire = cfg.PasswordExpire
}
