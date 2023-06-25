package courier

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
	"message-service/helpers"
	"message-service/types"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type EmailClient struct {
	username string
	password string
	host     string
	port     int
}

func InitEmailClient() EmailClient {
	convertPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	helpers.CheckForError(err, "Email port parsing failed")

	client := EmailClient{
		username: os.Getenv("SENDER_MAIL"),
		password: os.Getenv("SENDER_PASS"),
		host:     os.Getenv("SMTP_HOST"),
		port:     convertPort,
	}

	return client
}

func (client EmailClient) SendMessage(content types.Content) error {
	d := gomail.NewDialer(client.host, client.port, client.username, client.password)
	message := formatMessage(content, client.username)

	// Removing TLS verification on non-prod environments
	if os.Getenv("ENV") != "prod" {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if err := d.DialAndSend(message); err != nil {
		return err
	}

	log.Info("Email sent successfully")
	return nil
}

func formatMessage(content types.Content, from string) *gomail.Message {
	msg := gomail.NewMessage()

	addresses := make([]string, len(content.Recipients))
	for idx, recipient := range content.Recipients {
		addresses[idx] = msg.FormatAddress(recipient.Contact, recipient.Name)
	}

	msg.SetHeader("From", from)
	msg.SetHeader("To", addresses...)
	msg.SetHeader("Subject", string(content.Subject))
	msg.SetBody(content.Format, string(content.Content))

	return msg
}
