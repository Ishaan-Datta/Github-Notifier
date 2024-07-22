package main

type PullRequestResponse struct {
	Number     int `json:"number"`
	Repository struct {
		Name string `json:"name"`
		URL  string `json:"html_url"`
	} `json:"repository"`
	PullRequest struct {
		NodeID string `json:"node_id"`
		User   struct {
			Name string `json:"login"`
		} `json:"user"`
	} `json:"pull_request"`
}

type PushResponse struct {
	Repository struct {
		Name string `json:"name"`
		URL  string `json:"html_url"`
	} `json:"repository"`
	Pusher struct {
		Name string `json:"name"`
	} `json:"pusher"`
	HeadCommit struct {
		URL    string `json:"url"`
		Author struct {
			Name string `json:"name"`
		} `json:"author"`
		Message  string   `json:"message"`
		Time     string   `json:"timestamp"`
		Added    []string `json:"added"`
		Removed  []string `json:"removed"`
		Modified []string `json:"modified"`
	} `json:"head_commit"`

	Ref string `json:"ref"`
}

type IssueResponse struct {
	Repository struct {
		Name string `json:"name"`
		URL  string `json:"html_url"`
	} `json:"repository"`
	Issue struct {
		NodeID string `json:"node_id"`
		User   struct {
			Name string `json:"login"`
		} `json:"user"`
	} `json:"issue"`
}

type CreateResponse struct {
	Repository struct {
		Name string `json:"name"`
		URL  string `json:"html_url"`
	} `json:"repository"`
	Ref     string `json:"ref"`
	RefType string `json:"ref_type"`
	Master  string `json:"master_branch"`
}

type DeleteResponse struct {
	Repository struct {
		Name string `json:"name"`
		URL  string `json:"html_url"`
	} `json:"repository"`
	Ref     string `json:"ref"`
	RefType string `json:"ref_type"`
	Master  string `json:"master_branch"`
}
