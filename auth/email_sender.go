package auth

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

func EmailSender(receiverEmail string, subject string, body string, htmlBody string) error {
	emailApPassword := os.Getenv("email-pass")
	yourMail := os.Getenv("email")
	hostAddress := os.Getenv("hostaddress")
	emailName := os.Getenv("emailname")
	smtpAuthAddress := os.Getenv("emailauthaddress")
	// Create a new Email instance
	emailInstance := email.NewEmail()
	emailInstance.From = fmt.Sprintf("%s <%s>", emailName, yourMail)
	emailInstance.To = []string{receiverEmail}
	emailInstance.Subject = subject
	emailInstance.Text = []byte(body)
	emailInstance.HTML = []byte(htmlBody)

	// Setup TLS configuration for the SMTP server
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpAuthAddress,
	}
	// Send the email using TLS
	err := emailInstance.SendWithStartTLS(
		hostAddress,
		smtp.PlainAuth("", yourMail, emailApPassword, smtpAuthAddress),
		tlsConfig,
	)
	if err != nil {
		fmt.Println("There was an error sending the mail:", err)
		return err
	}
	return nil
}

func EmailVerification(receiverEmail string, code string) error {
	subject := "Verify your email"
	htmlBody := "<div> Your email verification code is: <span style=\"font-weight: bold;\">" + code + "</span></div>"
	err := EmailSender(receiverEmail, subject, "", htmlBody)
	if err != nil {
		return err
	}
	return nil
}
