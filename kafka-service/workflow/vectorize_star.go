package workflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type SyncUserStar struct {
	RepoId int64  `json:"repo_id"`
	Email  string `json:"email"`
}

func VectorizeStar(payload *SyncUserStar) (int, error) {
	token := os.Getenv("AUTHENTICATION_TOKEN")
	url := os.Getenv("TRANSFORMER_URL")

	jsonstring, err := json.Marshal(payload)
	if err != nil {
		return 400, err
	}

	req, err := http.NewRequest("POST", url+"/vectorize", bytes.NewBuffer(jsonstring))
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
