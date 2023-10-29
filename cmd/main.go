package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/ccallazans/match-notification-lambda/internal/configs"
	"github.com/ccallazans/match-notification-lambda/internal/service"
)

type Event struct {
	ID     uint   `json:"id"`
	Type   string `json:"type"`
	Topics []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"topics"`
	Message string `json:"message"`
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	log.Println("LAMBDA START")
	log.Printf("RECEIVED REQUEST: %+v\n", sqsEvent)

	cfg := configs.NewAWSConfig()
	sesClient := sesv2.NewFromConfig(cfg)
	sqsMessage := sqsEvent.Records[0]

	var event Event
	err := json.Unmarshal([]byte(sqsMessage.Body), &event)
	if err != nil {
		log.Printf("Error when parsing sqs message %s", err.Error())
		return err
	}

	var topics []string
	for _, topic := range event.Topics {
		topics = append(topics, topic.Name)
	}

	emails, err := service.GetUsersByTopic(topics)
	if err != nil {
		log.Printf("Error when geting user by topic %s", err.Error())
		return err
	}

	log.Printf("SENDING NOTIFICATION TO USERS: %+v\n", emails)

	subject := fmt.Sprintf("%s: %s", event.Type, event.Type)
	email := service.CreateEmail(emails, subject, event.Message)

	_, err = sesClient.SendEmail(context.TODO(), email)
	if err != nil {
		log.Printf("Error when sending email %s", err)
		return err
	}

	log.Println("END")

	return nil
}
