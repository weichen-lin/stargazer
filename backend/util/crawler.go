package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/weichen-lin/stargazer/domain"
)

var ErrNoToken = fmt.Errorf("invalid github token")

func GetUserStarredRepos(page int, token string) ([]domain.GithubRepository, error) {

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
		return nil, ErrNoToken
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var repos []domain.GithubRepository

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
