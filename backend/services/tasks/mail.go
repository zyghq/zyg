package tasks

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"

	"github.com/resend/resend-go/v2"
	"github.com/zyghq/zyg"
)

type KycMailData struct {
	PreviewText string
	MagicLink   string
}

func SendKycMail(to string, body string, verifyLink string) error {
	subject := "Started a new chat on Zyg."
	htmlTempl, err := template.ParseFiles("static/templates/mails/kyc.html")
	if err != nil {
		slog.Error("error parsing html template file", slog.Any("err", err))
		return err
	}
	// Read text template for non html text content.
	textTempl, err := template.ParseFiles("static/templates/mails/text/kyc.txt")
	if err != nil {
		slog.Error("error parsing text template file", slog.Any("err", err))
		return err
	}

	data := KycMailData{
		PreviewText: body,
		MagicLink:   verifyLink,
	}

	var htmlTemplOutput bytes.Buffer
	err = htmlTempl.Execute(&htmlTemplOutput, data)
	if err != nil {
		slog.Error("error executing html template", slog.Any("err", err))
		return err
	}
	htmlOutput := htmlTemplOutput.String()

	var textTemplOutput bytes.Buffer
	err = textTempl.Execute(&textTemplOutput, data)
	if err != nil {
		slog.Error("error executing text template", slog.Any("err", err))
		return err
	}
	textOutput := textTemplOutput.String()

	fmt.Println("************* HTML FOR MAIL **************")
	fmt.Println(htmlOutput)
	fmt.Println("************* END HTML FOR MAIL **************")

	fmt.Println("******* send mail to **********")
	fmt.Println(to)
	fmt.Println(subject)
	fmt.Println(textOutput)
	fmt.Println("******* END send mail to **********")

	apiKey := zyg.ResendApiKey()
	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "Sanchit <sanchit@updates.zyg.ai>",
		To:      []string{to},
		Subject: subject,
		Html:    htmlOutput,
		Text:    textOutput,
		ReplyTo: "sanchit@zyg.ai",
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		slog.Error("failed to send email", slog.Any("err", err))
		return err
	}

	slog.Info("sent email", slog.Any("Id", sent.Id))

	return nil
}
