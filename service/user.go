package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/weichen-lin/stargazer/db"
	"github.com/weichen-lin/stargazer/domain"
	"github.com/weichen-lin/stargazer/infrastructure"
)

func (s *Service) GetUser(ctx context.Context) (*domain.User, error) {
	var user *domain.User

	err := s.db.Run(ctx, func(q *db.Queries) error {
		var err error
		user, err = s.repository.GetUserByClerkId(ctx, q)
		if err != nil {
			return infrastructure.NewServiceError(
				ctx,
				http.StatusUnauthorized,
				"failed to get user",
				err.Error(),
			)
		}

		if user == nil {
			return infrastructure.NewServiceError(
				ctx,
				http.StatusUnauthorized,
				"failed to get user",
				"user not found",
			)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

// a golang map
var userMap = map[string]string{}

type OAuthAccessToken struct {
	Token string `json:"token"`
}

type ResponseBody []OAuthAccessToken

func (s *Service) GetOauthToken(ctx context.Context) (string, error) {
	clerkId, ok := infrastructure.GetClerkId(ctx)
	if !ok {
		return "", infrastructure.NewServiceError(
			ctx,
			http.StatusInternalServerError,
			"failed to get clerk id from context at GetOauthToken",
			"failed to get clerk id from context at GetOauthToken",
		)
	}

	if token, ok := userMap[clerkId]; ok {
		return token, nil
	}

	url := fmt.Sprintf("https://api.clerk.com/v1/users/%s/oauth_access_tokens/github", clerkId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", infrastructure.NewServiceError(
			ctx,
			http.StatusInternalServerError,
			"failed to create request",
			err.Error(),
		)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.ClerkSecret))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", infrastructure.NewServiceError(
			ctx,
			http.StatusInternalServerError,
			"failed to send request",
			err.Error(),
		)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", infrastructure.NewServiceError(
			ctx,
			http.StatusInternalServerError,
			"failed to get oauth token",
			fmt.Sprintf("failed to get oauth token, status code: %d", resp.StatusCode),
		)
	}

	var responseBody ResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return "", infrastructure.NewServiceError(
			ctx,
			http.StatusInternalServerError,
			"failed to decode response body",
			err.Error(),
		)
	}
	if len(responseBody) == 0 {
		return "", infrastructure.NewServiceError(
			ctx,
			http.StatusInternalServerError,
			"failed to get oauth token",
			"no oauth token found",
		)
	}
	token := responseBody[0].Token
	userMap[clerkId] = token

	return token, nil
}

type UserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type ClerkUser struct {
	Username       string         `json:"username"`
	EmailAddresses []EmailAddress `json:"email_addresses"`
}

type EmailAddress struct {
	EmailAddress string `json:"email_address"`
}

func (s *Service) GetUserInfo(ctx context.Context) (*UserInfo, error) {
	clerkId, ok := infrastructure.GetClerkId(ctx)
	if !ok {
		return nil, infrastructure.NewServiceError(
			ctx,
			http.StatusInternalServerError,
			"failed to get clerk id from context at GetUserInfo",
			"failed to get clerk id from context at GetUserInfo",
		)
	}

	url := fmt.Sprintf("https://api.clerk.com/v1/users/%s", clerkId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, infrastructure.NewServiceError(
			ctx,
			http.StatusInternalServerError,
			"failed to create request to get user info",
			err.Error(),
		)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.ClerkSecret))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, infrastructure.NewServiceError(
			ctx,
			http.StatusInternalServerError,
			"failed to send request to get user info",
			err.Error(),
		)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, infrastructure.NewServiceError(
			ctx,
			http.StatusInternalServerError,
			"failed to get user info",
			fmt.Sprintf("failed to get user info, status code: %d", resp.StatusCode),
		)
	}

	var responseBody ClerkUser
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return nil, infrastructure.NewServiceError(
			ctx,
			http.StatusInternalServerError,
			"failed to decode response body",
			err.Error(),
		)
	}
	if len(responseBody.EmailAddresses) == 0 {
		return nil, infrastructure.NewServiceError(
			ctx,
			http.StatusInternalServerError,
			"failed to get user info",
			"no email address found",
		)
	}
	email := responseBody.EmailAddresses[0].EmailAddress

	return &UserInfo{
		Email: email,
		Name:  responseBody.Username,
	}, nil
}

func (s *Service) UpdateUserCrontab(ctx context.Context, repository_count int32) error {
	return s.db.Run(ctx, func(q *db.Queries) error {
		return s.repository.SaveUserCrontab(ctx, q, repository_count)
	})
}
