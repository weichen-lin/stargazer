package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/weichen-lin/kafka-service/db"
)

func GetUserStarredRepos(page int, token string) ([]db.Repository, error) {

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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error code at get user stars: %d", resp.StatusCode)
	}

	var repos []db.Repository

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
