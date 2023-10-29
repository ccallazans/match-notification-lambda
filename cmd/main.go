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

	sqsMessage := sqsEvent.Records[0]
	event, err := parseSQSMessage(sqsMessage.Body)
	if err != nil {
		log.Printf("Error when parsing sqs message body %s", err.Error())
		return err
	}

	var topics []string
	for _, topic := range event.Topics {
		topics = append(topics, topic.Name)
	}

	emails, err := service.GetSubscriptionEmailsByTopic(topics)
	if err != nil {
		log.Printf("Error when geting user by topic %s", err.Error())
		return err
	}

	log.Printf("SENDING NOTIFICATION TO USERS: %+v\n", emails)
	err = sendEmail(event, emails)
	if err != nil {
		return err
	}

	log.Println("END")
	return nil
}

func parseSQSMessage(message string) (Event, error) {

	var event Event
	err := json.Unmarshal([]byte(message), &event)
	if err != nil {
		log.Printf("Error when parsing sqs message %s", err.Error())
		return Event{}, err
	}

	return event, nil
}

func sendEmail(event Event, emails []string) error {
	cfg := configs.NewAWSConfig()
	sesClient := sesv2.NewFromConfig(cfg)

	subject := fmt.Sprintf("MATCH-NOTIFICATION-API: %s", event.Type)
	email := service.CreateEmail(emails, subject, event.Message)

	_, err := sesClient.SendEmail(context.TODO(), email)
	if err != nil {
		log.Printf("Error when sending email %s", err)
		return err
	}

	return nil
}
