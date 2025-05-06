package util

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/resend/resend-go/v2"
)

type SendMailParams struct {
	Email      string
	Name       string
	StarsCount int
}

func SendMail(params *SendMailParams, apiKey string) error {
	templateFile := "templates/stars_finish.html"
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	templateData := struct {
		UserName   string
		StarsCount int
	}{
		UserName:   params.Name,
		StarsCount: params.StarsCount,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, templateData); err != nil {
		return err
	}

	client := resend.NewClient(apiKey)

	p := &resend.SendEmailRequest{
		From:    "Stargazer <stargazer@wei-chen.dev>",
		To:      []string{params.Email},
		Subject: "Your stars are ready!",
		Html:    body.String(),
	}

	_, err = client.Emails.Send(p)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return nil
}
