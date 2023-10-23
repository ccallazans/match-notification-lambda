package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/ccallazans/match-notification-lambda/internal/service"
)

type Event struct {
	ID      uint   `json:"id"`
	Topic   string `json:"topic"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	log.Println("LAMBDA START")
	log.Printf("RECEIVED REQUEST: %+v\n", sqsEvent)

	sqsMessage := sqsEvent.Records[0]

	event := &Event{}
	err := json.Unmarshal([]byte(sqsMessage.Body), &event)
	if err != nil {
		log.Printf("Error when parsing sqs message %s", err.Error())
		return err
	}

	users, err := service.GetUsersByTopic(event.Topic)
	if err != nil {
		log.Printf("Error when geting user by topic %s", err.Error())
		return err
	}

	subject := fmt.Sprintf("%s: %s", event.Type, event.Topic)
	log.Printf("SENDING NOTIFICATION TO USERS: %+v\n", users)

	input := &sesv2.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{"ccallazans@gmail.com"},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: &subject,
				},
				Body: &types.Body{
					Text: &types.Content{
						Data: &event.Message,
					},
				},
			},
		},
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Error creating aws config: %s", cfg)
	}

	sesClient := sesv2.NewFromConfig(cfg)
	result, err := sesClient.SendEmail(context.TODO(), input)
	if err != nil {
		log.Printf("Error when sending email %s", err)
		return err
	}

	log.Printf("Success!!: %s", result)
	log.Println("END")

	return nil
}
