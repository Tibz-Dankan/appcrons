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
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Email struct {
	Recipient string
}

// sends email to mails that match "@gmail.com" using smtp
func (e *Email) sendGmailHTML(html, subject string) error {
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

// sends email to mails of all categories using sendgrid-go client
func (e *Email) sendMailHTML(html, subject string) error {
	sendgridSenderMail := os.Getenv("SENDGRID_SENDER_MAIL")
	sendgridAPIKey := os.Getenv("SENDGRID_API_KEY")

	from := mail.NewEmail("Appcrons", sendgridSenderMail)
	to := mail.NewEmail("", e.Recipient)
	plainTextContent := ""
	htmlContent := html
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(sendgridAPIKey)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("response.StatusCode:", response.StatusCode)
	log.Println("response.Body:", response.Body)
	log.Println("response.Headers:", response.Headers)

	log.Println("Mail sent!")

	return nil
}

func (e *Email) isGmail(email string) bool {
	regex := regexp.MustCompile(`^[\w\.-]+@gmail\.com$`)
	return regex.MatchString(email)
}

func (u *Email) send(html, subject string) error {
	isGmail := u.isGmail(u.Recipient)

	if isGmail {
		if err := u.sendGmailHTML(html, subject); err != nil {
			log.Println(err)
			return err
		}
		return nil
	}

	if err := u.sendMailHTML(html, subject); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (e *Email) SendResetPassword(name, URL, subject string) error {
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
