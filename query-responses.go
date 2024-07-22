package main

import (
	"fmt"
	"reflect"
	"strings"
)

func iterateIssueFields(result IssueQueryResponse) IssueQueryResponse {
	refInterface := reflect.ValueOf(&result)
	nodeField := reflect.ValueOf(result).FieldByName("Node")

	// Iterate over fields in the Node struct
	for i := 0; i < nodeField.NumField(); i++ {
		fieldName := nodeField.Type().Field(i).Name
		fieldValue := nodeField.Field(i)

		if fieldValue.Kind() == reflect.String {
			if fieldValue.String() == "" {
				refInterface.Elem().FieldByName("Node").FieldByName(fieldName).SetString("(none)")
			}
		}

		if fieldValue.Kind() == reflect.Struct {
			nestedValue := reflect.ValueOf(fieldValue.Interface()).Field(0)
			if nestedValue.Len() == 0 {
				if fieldName == "Assignees" {
					value := reflect.ValueOf(struct {
						Nodes []struct {
							Login string
						}
					}{
						Nodes: []struct {
							Login string
						}{{Login: "(none)"}},
					})
					refInterface.Elem().FieldByName("Node").FieldByName(fieldName).Set(value)

				} else if fieldName == "Labels" {
					value := reflect.ValueOf(struct {
						Nodes []struct {
							Name string
						}
					}{
						Nodes: []struct {
							Name string
						}{{Name: "(none)"}},
					})
					refInterface.Elem().FieldByName("Node").FieldByName(fieldName).Set(value)

				} else if fieldName == "Author" {
					value := reflect.ValueOf(struct {
						Login string
					}{
						Login: "(none)",
					})
					refInterface.Elem().FieldByName("Node").FieldByName(fieldName).Set(value)
				}
			}
		}

	}
	return result
}

func iteratePRFields(result PullRequestQueryResponse) PullRequestQueryResponse {
	refInterface := reflect.ValueOf(&result)
	fmt.Printf("refInterface: %v\n", refInterface)
	nodeField := reflect.ValueOf(result).FieldByName("Node")

	// Iterate over fields in the Node struct
	for i := 0; i < nodeField.NumField(); i++ {
		fieldName := nodeField.Type().Field(i).Name
		fieldValue := nodeField.Field(i)

		if fieldValue.Kind() == reflect.Struct {
			nestedValue := reflect.ValueOf(fieldValue.Interface()).Field(0)
			if nestedValue.Len() == 0 {
				if fieldName == "Author" {
					value := reflect.ValueOf(struct {
						Login string
					}{
						Login: "(none)",
					})
					refInterface.Elem().FieldByName("Node").FieldByName(fieldName).Set(value)
				}
			}
		}
	}
	return result
}

func constructIssueMessage(event IssueQueryResponse, event2 IssueResponse) string {
	var labels []string
	for _, label := range event.Node.Labels.Nodes {
		labels = append(labels, label.Name)
	}

	var assignees []string
	for _, assignee := range event.Node.Assignees.Nodes {
		assignees = append(assignees, assignee.Login)
	}

	return fmt.Sprintf("Issue #%d: %s from repository: %s was updated to state: %s by %s at time: %s with the following body: %s and the following labels: %s, and the following assignees: %s. You can find the URL here: %s\n",
		event.Node.Number, event.Node.Title, event2.Repository.Name, event.Node.State, event.Node.Author.Login, event.Node.UpdatedAt, event.Node.BodyText, strings.Join(labels[:], ", "), strings.Join(assignees[:], ", "), event.Node.URL)
}

func constructPullRequestMessage(event PullRequestQueryResponse, event2 PullRequestResponse) string {
	return fmt.Sprintf("Pull request #%d: %s from repository: %s was updated to state: %s by %s at time: %s with the following changes: %d made on branch: %s. You can find the checks URL here: %s and the permalink here: %s\n",
		event.Node.Number, event.Node.Title, event2.Repository.Name, event.Node.State, event.Node.Author.Login, event.Node.UpdatedAt, event.Node.ChangedFiles, event.Node.BaseRefName, event.Node.ChecksURL, event.Node.Permalink)
}

func constructPushMessage(event PushResponse) string {
	return fmt.Sprintf("A push was made to branch: %s by %s with the head commit: %s at time: %s by %s adding files: %s, modifying files: %s, and removing files: %s. You can find the URL here: %s\n",
		event.Ref, event.Pusher.Name, event.HeadCommit.Message, event.HeadCommit.Time, event.HeadCommit.Author.Name, strings.Join(event.HeadCommit.Added[:], ", "), strings.Join(event.HeadCommit.Modified[:], ", "), strings.Join(event.HeadCommit.Removed[:], ", "), event.HeadCommit.URL)
}

func constructDeleteMessage(event DeleteResponse) string {
	return fmt.Sprintf("%s %s deleted on repository %s (%s) with master branch: %s", event.RefType, event.Ref, event.Repository.Name, event.Repository.URL, event.Master)
}

func constructCreateMessage(event CreateResponse) string {
	return fmt.Sprintf("%s %s created on repository %s (%s) with master branch: %s", event.RefType, event.Ref, event.Repository.Name, event.Repository.URL, event.Master)
}
