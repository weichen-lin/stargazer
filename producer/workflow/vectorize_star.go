package workflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type TransFormerParam struct {
	RepoId int64  `json:"repo_id"`
	Name   string `json:"name"`
}

type SyncUserStarMsg struct {
	RepoId   int64  `json:"repo_id"`
	UserName string `json:"username"`
}

func VectorizeStar(payload *SyncUserStarMsg) (int, error) {
	token := os.Getenv("AUTHENTICATION_TOKEN")

	url := "https://stargazer-transformer-qmy6az4wfa-de.a.run.app/vectorize"

	params := TransFormerParam{
		RepoId: payload.RepoId,
		Name:   payload.UserName,
	}

	jsonstring, err := json.Marshal(params)
	if err != nil {
		return 400, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonstring))
	if err != nil {
		return 400, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 400, err
	}
	defer resp.Body.Close()

	// Check the response message
	body := new(bytes.Buffer)
	body.ReadFrom(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return 400, fmt.Errorf("error: %d", resp.StatusCode)
	}

	return resp.StatusCode, nil
}
