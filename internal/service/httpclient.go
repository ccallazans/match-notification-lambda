package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ccallazans/match-notification-lambda/internal/models"
)

func GetUsersByTopic(topicName string) ([]models.User, error) {
	apiUrl := os.Getenv("MATCH_NOTIFICATION_API")

	response, err := http.Get(fmt.Sprintf("%s/users?topic=%s", apiUrl, topicName))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading response body: %s", err)
		return nil, err
	}

	users := []models.User{}
	if err := json.Unmarshal(responseBody, &users); err != nil {
		log.Printf("Error parsing JSON: %s", err)
		return nil, err
	}

	return users, nil
}