package workflow

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Owner struct {
	AvatarURL string `json:"avatar_url"`
}

type Repository struct {
	ID              int       `json:"id"`
	FullName        string    `json:"full_name"`
	Owner           Owner     `json:"owner"`
	HTMLURL         string    `json:"html_url"`
	Description     string    `json:"description"`
	UpdatedAt       time.Time `json:"updated_at"`
	StargazersCount int       `json:"stargazers_count"`
	Language        string    `json:"language"`
	DefaultBranch   string    `json:"default_branch"`
}

func GetUserStarredRepos() ([]byte, error) {

	req, err := http.NewRequest("GET", "https://api.github.com/users/weichen-lin/starred", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+"")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repos []Repository

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &repos)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	fmt.Println("Starred Repositories:", repos)

	return body, nil
}
