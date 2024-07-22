package main

import (
	"fmt"
	"log"
	"net/http"
)

// Task represents a unit of work to be performed
type Task struct {
	message string
}

// TaskResult holds the result of a task
type TaskResult struct {
	application string
	status      bool
	reason      string
}

var taskQueue = make(chan Task, 10)
var resultQueue = make(chan TaskResult, 10)
var config, err = readConfig()

// processes tasks from the taskQueue
func Worker(id int) {
	for task := range taskQueue {
		fmt.Printf("Dispatching task with message: %s", task.message)
		go sendDiscordMessage(DiscordPayload{Content: task.message})
		go sendSlackMessage(task.message)
	}
}

// processes tasks from the resultQueue
func Reporter() {
	for result := range resultQueue {
		if result.status {
			fmt.Printf("Task completed successfully by %s\n", result.application)
		} else {
			fmt.Printf("Task failed by %s, with reason: %s\n", result.application, result.reason)
		}
	}
}

func main() {
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	http.HandleFunc("/webhook", webhookHandler)
	fmt.Println("Starting server on :3000...")

	numWorkers := 3

	for i := 0; i < numWorkers; i++ {
		go Worker(i + 1)
	}

	go Reporter()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
