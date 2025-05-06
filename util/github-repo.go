package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type GithubRepository struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Owner       Owner    `json:"owner"`
	HTMLURL     string   `json:"html_url"`
	Description string   `json:"description"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	Homepage    string   `json:"homepage"`
	Stargazers  int      `json:"stargazers_count"`
	Language    string   `json:"language"`
	Archived    bool     `json:"archived"`
	Topics      []string `json:"topics"`
	Forks       int      `json:"forks"`
	OpenIssues  int      `json:"open_issues"`
	Watchers    int      `json:"watchers"`
}

type Owner struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
}

func GetUserStarredRepos(page int, token string) ([]GithubRepository, error) {

	url := fmt.Sprintf("https://api.github.com/user/starred?&page=%d", page)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("unauthorized: invalid token")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var repos []GithubRepository

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &repos)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %s", err)
	}

	return repos, nil
}
