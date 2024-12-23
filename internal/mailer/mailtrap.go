package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"time"

	gomail "gopkg.in/mail.v2"
)

type mailtrapClient struct {
	fromEmail string
	apiKey    string
}

func NewMailTrapClient(apiKey, fromEmail string) (mailtrapClient, error) {
	if apiKey == "" {
		return mailtrapClient{}, errors.New("api key is required")
	}

	return mailtrapClient{
		fromEmail: fromEmail,
		apiKey:    apiKey,
	}, nil
}

func (m mailtrapClient) Send(templateFile, username, email string, data any, isSandBox bool) (int, error) {
	// template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
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
		return -1, err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject.String())
	message.AddAlternative("text/mail", body.String())

	dailer := gomail.NewDialer("sandbox.smtp.mailtrap.io", 587, "d4c13dd2a25044", m.apiKey)

	var retryErr error
	for i := 0; i < maxRetries; i++ {
		if err := dailer.DialAndSend(message); err != nil {
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return 200, nil
	}
	return -1, fmt.Errorf("failed to send email after %d attempt, error: %v", maxRetries, retryErr)
}
