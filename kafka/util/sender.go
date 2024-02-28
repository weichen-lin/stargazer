package util

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
)

type SendMailParams struct {
	Email     string
	Name 	string
	StarsCount int
}

func SendMail(params *SendMailParams) error {

	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	senderEmail := "gitstargazer@gmail.com"
	senderPassword := os.Getenv("APP_PASSWORD")

	templateFile := "templates/stars_finish.html"
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	templateData := struct {
		UserName string
		StarsCount int
	} {
		UserName: params.Name,
		StarsCount: params.StarsCount,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, templateData); err != nil {
		return err
	}

	msg := []byte(
		"From: " + "StarGazer <gitstargazer@gmail.com>" + "\r\n" +
			"To: " + params.Email + "\r\n" +
			"Subject: Completion Procedure for StarGazer\r\n" +
			"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
			"\r\n" + body.String())

	err = smtp.SendMail(smtpHost+":"+fmt.Sprint(smtpPort), auth, senderEmail, []string{params.Email}, msg)
	if err != nil {
		return err
	}

	return nil
}