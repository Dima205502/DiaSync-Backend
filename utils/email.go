package utils

import (
	"net/smtp"
)

func SendVerifyTokenMail(email, verifyEmailToken string) error {
	auth := smtp.PlainAuth("", sender, app_password, smtp_server)

	msg := "Subject: Verify Email\nClick on the link to confirm your email\nhttp://" +
		server_adr + "/auth/verify-email?token=" + verifyEmailToken

	err := smtp.SendMail(smtp_adr, auth, sender, []string{email}, []byte(msg))

	return err
}

func SendNewPasswordEmail(email, newPasswordToken string) error {
	auth := smtp.PlainAuth("", sender, app_password, smtp_server)

	msg := "Subject: Verify Email\nClick on the link to confirm your email\nhttp://" +
		server_adr + "/auth/verify-newpassword?token=" + newPasswordToken

	err := smtp.SendMail(smtp_adr, auth, sender, []string{email}, []byte(msg))

	return err
}
