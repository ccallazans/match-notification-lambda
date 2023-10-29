package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type User struct {
	ID     uint   `json:"id"`
	Email  string `json:"email"`
	Topics []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"topics"`
}

func GetUsersByTopic(topics []string) ([]string, error) {
	apiUrl := os.Getenv("MATCH_NOTIFICATION_API")

	var queryString string
	for _, topic := range topics {
		queryString += fmt.Sprintf("?topic=%s", topic)
	}

	response, err := http.Get(fmt.Sprintf("%s/api/v1/subscriptions%s", apiUrl, queryString))
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

	users := []User{}
	if err := json.Unmarshal(responseBody, &users); err != nil {
		log.Printf("Error parsing JSON: %s", err)
		return nil, err
	}

	emails := make(map[string]bool)
	for _, user := range users {
		emails[user.Email] = true
	}

	keys := make([]string, 0, len(emails))
	for k := range emails {
		keys = append(keys, k)
	}

	return keys, nil
}
