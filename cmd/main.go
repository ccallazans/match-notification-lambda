package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ccallazans/match-notification-lambda/internal/service"
)

type Event struct {
	ID      string `json:"id"`
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

	log.Printf("SENDING NOTIFICATION TO USERS: %+v\n", users)
	log.Println("END")

	return nil
}
