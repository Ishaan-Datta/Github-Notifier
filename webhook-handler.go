package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	eventType := r.Header.Get("X-GitHub-Event")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	responseStruct, eventName := GetEventDetails(eventType, body)
	if responseStruct == nil {
		http.Error(w, "Failed to process event", http.StatusBadRequest)
		return
	}

	determineQueryType(responseStruct, eventName)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received"))
}

func GetEventDetails(event string, body []byte) (interface{}, string) {
	if strings.Contains(event, "push") {
		fmt.Printf("Received a push event\n")
		var response PushResponse
		fillStruct(body, &response)
		return response, "push"
	} else if strings.Contains(event, "issue") {
		fmt.Printf("Received an issue event\n")
		var response IssueResponse
		fillStruct(body, &response)
		return response, "issue"
	} else if strings.Contains(event, "pull_request") {
		fmt.Printf("Received a pull request event\n")
		var response PullRequestResponse
		fillStruct(body, &response)
		return response, "pull request"
	} else if strings.Contains(event, "create") {
		fmt.Printf("Received a create event\n")
		var response CreateResponse
		fillStruct(body, &response)
		return response, "create"
	} else if strings.Contains(event, "delete") {
		fmt.Printf("Received a delete event\n")
		var response DeleteResponse
		fillStruct(body, &response)
		return response, "delete"
	} else {
		return nil, "nil"
	}
}

func fillStruct(body []byte, response interface{}) error {
	err := json.Unmarshal(body, response)
	return err
}
