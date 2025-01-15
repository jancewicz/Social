package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"time"

	gomail "gopkg.in/gomail.v2"
)

type mailtrapClient struct {
	fromEmail string
	apiKey    string
}

func NewMailTrapClient(apiKey, fromEmail string) (mailtrapClient, error) {
	if apiKey == "" {
		return mailtrapClient{}, errors.New("api key not provided")
	}

	return mailtrapClient{
		fromEmail: fromEmail,
		apiKey:    apiKey,
	}, nil
}

func (m mailtrapClient) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	// Template building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	fmt.Println(tmpl)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return -1, nil
	}

	message := gomail.NewMessage()
	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject.String())

	message.AddAlternative("text/html", body.String())

	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", m.apiKey)

	if err := dialer.DialAndSend(message); err != nil {
		return -1, err
	}

	for i := 0; i < maxRetires; i++ {
		if err := dialer.DialAndSend(message); err != nil {
			log.Printf("Failed to send email to %v, attempt %d of %d", email, i+1, maxRetires)
			log.Printf("Error: %v", err.Error())

			// backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
	}

	return 200, nil
}
