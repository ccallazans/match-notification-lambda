package service

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

func CreateEmail(receivers []string, subject string, body string) *sesv2.SendEmailInput {
	sender := os.Getenv("FROM_EMAIL_ADDRESS")

	return &sesv2.SendEmailInput{
		FromEmailAddress: &sender,
		Destination: &types.Destination{
			ToAddresses: receivers,
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: &subject,
				},
				Body: &types.Body{
					Text: &types.Content{
						Data: &body,
					},
				},
			},
		},
	}
}
