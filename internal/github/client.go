package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	token   string
	host    string
	baseURL string
}

type PullRequest struct {
	Number     int    `json:"number"`
	Title      string `json:"title"`
	URL        string `json:"html_url"`
	Repository string `json:"-"`
}

type searchResponse struct {
	Items []PullRequest `json:"items"`
}

func NewClient(token, host string) *Client {
	baseURL := "https://api.github.com"
	if host != "" && host != "github.com" {
		// GitHub Enterprise
		baseURL = fmt.Sprintf("https://%s/api/v3", host)
	}

	return &Client{
		token:   token,
		host:    host,
		baseURL: baseURL,
	}
}

func (c *Client) GetReviewRequests(repo, username string) ([]PullRequest, error) {
	// Build search query
	query := fmt.Sprintf("type:pr state:open review-requested:%s repo:%s", username, repo)
	
	url := fmt.Sprintf("%s/search/issues?q=%s&per_page=100", c.baseURL, strings.ReplaceAll(query, " ", "+"))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	var searchResp searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Add repository name to each PR
	for i := range searchResp.Items {
		searchResp.Items[i].Repository = repo
	}

	return searchResp.Items, nil
}
