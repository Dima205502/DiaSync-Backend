package utils

import (
	"net/smtp"
)

func SendVerifyTokenMail(email, verifyEmailToken string) error {
	auth := smtp.PlainAuth("", sender, appPassword, smtpServer)

	msg := "Subject: Verify Email\nClick on the link to confirm your email\nhttp://" +
		serverAdr + "/auth/verify-email?token=" + verifyEmailToken

	err := smtp.SendMail(smtpAdr, auth, sender, []string{email}, []byte(msg))

	return err
}

func SendNewPasswordEmail(email, newPasswordToken string) error {
	auth := smtp.PlainAuth("", sender, appPassword, smtpServer)

	msg := "Subject: Verify Email\nClick on the link to confirm your email\nhttp://" +
		serverAdr + "/auth/verify-newpassword?token=" + newPasswordToken

	err := smtp.SendMail(smtpAdr, auth, sender, []string{email}, []byte(msg))

	return err
}
