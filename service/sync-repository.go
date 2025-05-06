package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/weichen-lin/stargazer/db"
	"github.com/weichen-lin/stargazer/infrastructure"
	"github.com/weichen-lin/stargazer/util"
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

func parseTime(t string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return time.Time{}, err
	}

	if parsedTime.Year() < 1970 {
		return time.Time{}, errors.New("year cannot be less than 1970")
	}

	return parsedTime.UTC(), nil
}

func (s *Service) SaveGithubRepository(ctx context.Context, githubRepository *util.GithubRepository) error {
	return s.db.RunWithTransaction(ctx, func(q *db.Queries) error {
		owner, err := q.UpsertOwner(context.Background(), db.UpsertOwnerParams{
			ID:        int32(githubRepository.Owner.ID),
			Name:      githubRepository.Owner.Login,
			AvatarUrl: githubRepository.Owner.AvatarURL,
		})
		if err != nil {
			return infrastructure.NewServiceError(ctx, http.StatusInternalServerError, "Failed to upsert owner", err.Error())
		}

		createdAt, err := parseTime(githubRepository.CreatedAt)
		if err != nil {
			return infrastructure.NewServiceError(ctx, http.StatusInternalServerError, "Failed to parse created_at", err.Error())
		}

		updatedAt, err := parseTime(githubRepository.UpdatedAt)
		if err != nil {
			return infrastructure.NewServiceError(ctx, http.StatusInternalServerError, "Failed to parse updated_at", err.Error())
		}

		repository, err := q.UpsertRepository(context.Background(), db.UpsertRepositoryParams{
			ID:          int32(githubRepository.ID),
			Name:        githubRepository.Name,
			OwnerID:     owner.ID,
			HtmlUrl:     githubRepository.HTMLURL,
			Homepage:    &githubRepository.Homepage,
			Description: &githubRepository.Description,
			Watchers:    int32(githubRepository.Watchers),
			Forks:       int32(githubRepository.Forks),
			OpenIssues:  int32(githubRepository.OpenIssues),
			Language:    githubRepository.Language,
			Archived:    githubRepository.Archived,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})
		if err != nil {
			return infrastructure.NewServiceError(ctx, http.StatusInternalServerError, "Failed to upsert repository", err.Error())
		}

		for _, topic := range githubRepository.Topics {
			topic, err := q.UpsertTopics(context.Background(), topic)
			if err != nil {
				return infrastructure.NewServiceError(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to upsert topic: %s", topic), err.Error())
			}

			_, err = q.UpsertRepositoryTopics(context.Background(), db.UpsertRepositoryTopicsParams{
				RepoID:  repository.ID,
				TopicID: topic.ID,
			})
			if err != nil {
				return infrastructure.NewServiceError(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to upsert repository topic: %s", topic), err.Error())
			}
		}

		return nil
	})
}

func (s *Service) SaveUserStarredRepository(ctx context.Context, repository_id int32) error {
	return s.db.Run(ctx, func(q *db.Queries) error {
		return s.repository.SaveUserStarredRepository(ctx, q, repository_id)
	})
}
