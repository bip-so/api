package pkg

import (
	"context"
	"fmt"
	configv2 "github.com/aws/aws-sdk-go-v2/config"
	credv2 "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
)

type BipMailer struct{}

type Recipient struct {
	ToEmails  []string
	CcEmails  []string
	BccEmails []string
}

func (m BipMailer) getSESClient() (client *ses.Client, err error) {
	cfg, err := configv2.LoadDefaultConfig(
		context.TODO(),
		configv2.WithDefaultRegion(configs.GetConfigString("S3_REGION")),
		configv2.WithCredentialsProvider(
			credv2.NewStaticCredentialsProvider(configs.GetConfigString("SES_ACCESS_KEY_ID"), configs.GetConfigString("SES_SECRET_ACCESS_KEY"), ""),
		))
	if err != nil {
		return client, fmt.Errorf("SES aws configuration. Error: %v", err)
	}
	client = ses.NewFromConfig(cfg)
	return client, err
}

// Send email via SES
func (m BipMailer) SendEmail(to []string, cc []string, bcc []string, subject string, bodyHTML string, bodyText string) error {
	var UTF8CHARSET = "UTF-8"
	// SES Client
	client, creatingClient := m.getSESClient()
	if creatingClient != nil {
		fmt.Printf("error while creating session for sending email %v", creatingClient)
		fmt.Println(creatingClient)
	}

	sender := configs.GetCurrentSystemEmail()

	// chirag@stage-emails.dev.bip.so
	//Sender = "santhosh@stage-emails.dev.bip.so"
	//Sender = "paras@stage-emails.dev.bip.so"
	// making the message !
	message := &types.Message{
		Body: &types.Body{
			Html: &types.Content{
				Data:    &bodyHTML,
				Charset: &UTF8CHARSET, // UTF-8, ISO-8859-1,
			},
			Text: &types.Content{
				Data: &bodyText,
				// UTF-8, ISO-8859-1,
				Charset: &UTF8CHARSET,
			},
		},
		Subject: &types.Content{
			Data:    &subject,
			Charset: &UTF8CHARSET,
		},
	}
	// Email Input
	emailInput := ses.SendEmailInput{
		Message: message,
		Destination: &types.Destination{
			ToAddresses:  to,
			CcAddresses:  cc,
			BccAddresses: bcc,
		},
		Source: &sender, // Match SES
	}
	result, err := client.SendEmail(context.TODO(), &emailInput)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("Email sent successfully. Response: %v", result)
	return nil
}

func TestMailer() error {

	body := "<html></head><title>This is html body</title></head><body><h>hello there</h><br>this is message body</body></html>"
	bodyPlainText := "This is message body"
	subject := "Test Emails"
	var mailer BipMailer
	toList := []string{"chirax@gmail.com", "santhoshkumar9713@gmail.com", "nitish.rddy@gmail.com"}
	emptyList := []string{}
	err := BipMailer.SendEmail(mailer, toList, emptyList, emptyList, subject, body, bodyPlainText)
	if err != nil {
		return err
	}

	return nil
}
