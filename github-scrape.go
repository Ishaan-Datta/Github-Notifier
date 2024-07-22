package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/machinebox/graphql"
)

type IssueQueryResponse struct {
	Node struct {
		Title     string
		Assignees struct {
			Nodes []struct {
				Login string
			}
		}
		Author struct {
			Login string
		}
		BodyText string
		Labels   struct {
			Nodes []struct {
				Name string
			}
		}
		Number    int
		URL       string
		State     string
		UpdatedAt string
	}
}

type PullRequestQueryResponse struct {
	Node struct {
		Title        string
		BaseRefName  string
		ChangedFiles int
		ChecksURL    string
		UpdatedAt    string
		Number       int
		Permalink    string
		State        string
		Author       struct {
			Login string
		}
	}
}

func determineQueryType(event interface{}, eventType string) {
	if eventType == "push" {
		MessageString := constructPushMessage(event.(PushResponse))
		// fmt.Printf("%s", MessageString)
		taskQueue <- Task{message: MessageString}
	} else if eventType == "pull request" {
		queryResult := QueryGitHubGraphQL(pullRequestQuery(event.(PullRequestResponse)))
		MessageString := constructPullRequestMessage(queryResult.(PullRequestQueryResponse), event.(PullRequestResponse))
		// fmt.Printf("%s", MessageString)
		taskQueue <- Task{message: MessageString}
	} else if eventType == "issue" {
		queryResult := QueryGitHubGraphQL(issueQuery(event.(IssueResponse)))
		MessageString := constructIssueMessage(queryResult.(IssueQueryResponse), event.(IssueResponse))
		// fmt.Printf("%s", MessageString)
		taskQueue <- Task{message: MessageString}
	} else if eventType == "create" {
		MessageString := constructCreateMessage(event.(CreateResponse))
		// fmt.Printf("%s", MessageString)
		taskQueue <- Task{message: MessageString}
	} else if eventType == "delete" {
		MessageString := constructDeleteMessage(event.(DeleteResponse))
		// fmt.Printf("%s", MessageString)
		taskQueue <- Task{message: MessageString}
	} else {
		fmt.Printf("Unknown event type %s:", eventType)
	}
}

func issueQuery(event IssueResponse) string {
	return fmt.Sprintf(`
		query {
			node(id: "%s") {
				... on Issue {
					title
					assignees(first: 10) {
						nodes {
							login
						}
					}
					author {
						login
					}
					body
					labels(first: 10) {
						nodes {
							name
						}
					}
					number
					url
					state
					updatedAt
				}
			}
		}`, string(event.Issue.NodeID))
}

func pullRequestQuery(event PullRequestResponse) string {
	return fmt.Sprintf(`
		query {
			node(id: "%s") {
				... on PullRequest {
					title
					baseRefName
					changedFiles
					checksUrl
					updatedAt
					number
					permalink
					state
					author {
						login
					}
				}	
			}
		}
	`, string(event.PullRequest.NodeID))
}

type CustomTransport struct {
	Transport http.RoundTripper
}

func (t *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.Transport.RoundTrip(req)
}

func QueryGitHubGraphQL(query string) interface{} {
	token := config.GithubOAuthToken

	client := graphql.NewClient("https://api.github.com/graphql", graphql.WithHTTPClient(&http.Client{
		Transport: &CustomTransport{
			Transport: http.DefaultTransport,
		},
	}))
	request := graphql.NewRequest(query)
	request.Header.Set("Authorization", "Bearer "+token)

	var rawResponse json.RawMessage
	err := client.Run(context.Background(), request, &rawResponse)
	if err != nil {
		fmt.Println("Error querying GitHub GraphQL API:", err)
		return nil
	}

	// fmt.Printf("Raw JSON Response: %s\n", rawResponse)

	var result interface{}
	if queryContains(query, "Issue") {
		var issueData IssueQueryResponse
		err = json.Unmarshal(rawResponse, &issueData)
		if err != nil {
			fmt.Println("Error unmarshalling Issue data:", err)
			return nil
		}
		result = iterateIssueFields(issueData)
	} else if queryContains(query, "PullRequest") {
		var pullRequestData PullRequestQueryResponse
		err = json.Unmarshal(rawResponse, &pullRequestData)
		if err != nil {
			fmt.Println("Error unmarshalling PullRequest data:", err)
			return nil
		}
		result = iteratePRFields(pullRequestData)
	}
	return result
}

func queryContains(query string, substr string) bool {
	return bytes.Contains([]byte(query), []byte(substr))
}
