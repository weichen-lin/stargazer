package workflow

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Owner struct {
	AvatarURL string `json:"avatar_url"`
}

type Repository struct {
	ID              int64  `json:"id"`
	FullName        string `json:"full_name"`
	Owner           Owner  `json:"owner"`
	HTMLURL         string `json:"html_url"`
	Description     string `json:"description"`
	UpdatedAt       string `json:"updated_at"`
	StargazersCount int    `json:"stargazers_count"`
	Language        string `json:"language"`
	DefaultBranch   string `json:"default_branch"`
}

type GetGithubReposInfo struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	Page     int    `json:"page"`
}

func GetUserStarredRepos(info *GetGithubReposInfo) ([]Repository, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN not set")
	}

	url := fmt.Sprintf("https://api.github.com/users/%s/starred?&page=%d", info.Username, info.Page)

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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %d", resp.StatusCode)
	}

	var repos []Repository

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &repos)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	return repos, nil
}
