package services

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	gmail "github.com/go-mail/mail"
	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type Email struct {
	Recipient string
}

// sends email to mails that match "@gmail.com" using smtp
func (e *Email) SendGmailHTML(html, subject string) error {
	senderMail := os.Getenv("GMAIL_SENDER_MAIL")
	appPassword := os.Getenv("GMAIL_APP_PASSWORD")

	m := gmail.NewMessage()
	m.SetHeader("From", senderMail)
	m.SetHeader("To", e.Recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", html)

	d := gmail.NewDialer("smtp.gmail.com", 587, senderMail, appPassword)

	d.StartTLSPolicy = gmail.MandatoryStartTLS

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	log.Println("Mail sent!")

	return nil
}

// sends email to mails of all categories using mailjet wrapper
func (e *Email) sendMailHTML(html, subject string) error {
	publicKey := os.Getenv("MJ_APIKEY_PUBLIC")
	secretKey := os.Getenv("MJ_APIKEY_PRIVATE")
	senderMail := os.Getenv("MJ_SENDER_MAIL")

	mailjetClient := mailjet.NewMailjetClient(publicKey, secretKey)
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: senderMail,
				Name:  "Appcrons",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: e.Recipient,
					Name:  "User x",
				},
			},
			Subject:  subject,
			TextPart: "",
			HTMLPart: html,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}

	if _, err := mailjetClient.SendMailV31(&messages); err != nil {
		return err
	}

	log.Println("mail sent!")

	return nil
}

func (e *Email) IsGmail(email string) bool {
	regex := regexp.MustCompile(`^[\w\.-]+@gmail\.com$`)
	return regex.MatchString(email)
}

func (u *Email) send(html, subject string) error {
	if err := u.sendMailHTML(html, subject); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (e *Email) SendResetPassword(name, URL, subject string) error {
	if os.Getenv("GO_ENV") == "testing" || os.Getenv("GO_ENV") == "staging" {
		return nil
	}
	data := struct {
		Subject string
		Name    string
		URL     string
		Year    string
	}{
		Subject: subject,
		Name:    name,
		URL:     URL,
		Year:    strconv.Itoa(time.Now().Year()),
	}

	var body bytes.Buffer

	templatePath, err := filepath.Abs("./internal/templates/email/reset-password.html")
	if err != nil {
		log.Println("Error finding absolute path:", err)
		return err
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Println("Error parsing template:", err)
		return err
	}

	err = tmpl.Execute(&body, data)
	if err != nil {
		log.Println("Error executing template:", err)
		return err
	}

	if err := e.send(body.String(), subject); err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	return nil
}
