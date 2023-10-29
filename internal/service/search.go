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

func GetSubscriptionEmailsByTopic(topics []string) ([]string, error) {
	apiUrl := os.Getenv("MATCH_NOTIFICATION_API")

	var queryString string
	for _, topic := range topics {
		queryString += fmt.Sprintf("&topic=%s", topic)
	}

	url := fmt.Sprintf("%s/api/v1/subscriptions?%s", apiUrl, queryString)
	
	log.Printf("Sending Get Request topics: %s to: %s", topics, url)
	response, err := http.Get(url)
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

	log.Printf("RECEIVED RESPONSE: %s\n", responseBody)
	emails, err := parseResponse(responseBody)
	if err != nil {
		return nil, err
	}

	return emails, nil
}

func parseResponse(response []byte) ([]string, error) {

	var users []User
	if err := json.Unmarshal(response, &users); err != nil {
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
