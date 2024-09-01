package util

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/weichen-lin/kabaka"
	"github.com/weichen-lin/stargazer/db"
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

type GetGithubReposInfo struct {
	Email string `json:"email"`
	Page  int    `json:"page"`
}

func GetGithubRepos(database *db.Database, msg kabaka.Message, writer *kabaka.Kabaka) error {

	var info GetGithubReposInfo

	err := json.Unmarshal(msg.Value, &info)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %s", err.Error())
	}

	ctx, err := db.WithEmail(context.Background(), info.Email)
	if err != nil {
		return err
	}

	user, err := database.GetUser(ctx)

	stars, err := GetUserStarredRepos(info.Page, user.AccessToken())

	for _, star := range stars {
		repo, err := domain.NewRepository(&star)
		if err != nil {
			return err
		}

		err = database.CreateRepository(ctx, repo)
		if err != nil {
			return err
		}
	}

	if len(stars) == 30 {
		info.Page++

		jsonString, err := json.Marshal(info)
		if err != nil {
			return fmt.Errorf("error marshalling JSON: %s", err.Error())
		}

		err = writer.Publish("star-syncer", jsonString)

		if err != nil {
			return fmt.Errorf("error sending message: %s", err.Error())
		}

	} else {
		starsCount := (info.Page-1)*30 + len(stars)

		SendMail(&SendMailParams{
			Email:      info.Email,
			Name:       user.Name(),
			StarsCount: starsCount,
		})
	}

	return nil
}
