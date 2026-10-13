package common

import (
	"crypto/tls"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

func SendMail(to string, subject string, body string) error {
	host := viper.GetString("mail.host")
	port := viper.GetInt("mail.port")
	username := viper.GetString("mail.username")
	password := viper.GetString("mail.password")
	from := viper.GetString("mail.from")

	if host == "" || host == "sandbox.smtp.mailtrap.io" && username == "your_mailtrap_username" {
		// Fallback to print if mailtrap not configured
		fmt.Printf("[MAIL MOCK] To: %s | Subject: %s | Body: %s\n", to, subject, body)
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func SendTemplateEmail(to string, tmpl EmailTemplate) error {
	host := viper.GetString("mail.host")
	port := viper.GetInt("mail.port")
	username := viper.GetString("mail.username")
	password := viper.GetString("mail.password")
	from := viper.GetString("mail.from")

	if host == "" || host == "sandbox.smtp.mailtrap.io" && username == "your_mailtrap_username" {
		fmt.Printf("[MAIL MOCK TEMPLATE] To: %s | Subject: %s\n", to, tmpl.Subject())
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", tmpl.Subject())
	m.SetBody("text/plain", tmpl.PlainTextBody())
	m.AddAlternative("text/html", tmpl.HtmlBody())

	d := gomail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send template email: %v", err)
	}

	return nil
}
