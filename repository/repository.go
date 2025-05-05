package repository

import (
	"context"

	"github.com/weichen-lin/stargazer/db"
	"github.com/weichen-lin/stargazer/domain"
)

type LanguageDistribution struct {
	Language string `json:"language"`
	Count    int64  `json:"count"`
}

func (r *Repository) GetLanguageDistribution(ctx context.Context, q *db.Queries, user *domain.User) ([]*LanguageDistribution, error) {
	languageDistributions, err := q.GetStarredLanguageDistributionByUser(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	languageDistributionList := make([]*LanguageDistribution, 0, len(languageDistributions))
	for _, languageDistribution := range languageDistributions {
		languageDistributionList = append(languageDistributionList, &LanguageDistribution{
			Language: languageDistribution.Language,
			Count:    languageDistribution.Count,
		})
	}

	return languageDistributionList, nil
}

func (r *Repository) GetLatestUserUpdatedRepositories(ctx context.Context, q *db.Queries, user *domain.User) ([]*domain.Repository, error) {
	repositories, err := q.ListStarredRepositoriesByUpdatedAt(ctx, db.ListStarredRepositoriesByUpdatedAtParams{
		UserID: user.Id,
		Limit:  10,
	})
	if err != nil {
		return nil, err
	}

	repositoryList := make([]*domain.Repository, 0, len(repositories))
	for _, repository := range repositories {

		repo, err := domain.FromRepositoryEntity(
			&domain.RepositoryEntity{
				Id:          repository.ID,
				Name:        repository.Name,
				OwnerName:   repository.OwnerName,
				AvatarURL:   repository.OwnerAvatarUrl,
				HtmlURL:     repository.HtmlUrl,
				Homepage:    *repository.Homepage,
				Description: *repository.Description,
				CreatedAt:   repository.CreatedAt,
				UpdatedAt:   repository.UpdatedAt,
				SyncedAt:    repository.SyncedAt,
				Watchers:    repository.Watchers,
				Forks: 	repository.Forks,
				OpenIssues:  repository.OpenIssues,
				Language:    repository.Language,
				Archived:    repository.Archived,
				Topics:      repository.Topics,
			},
		)

		if err != nil {
			return nil, err
		}

		repositoryList = append(repositoryList, repo)
	}

	return repositoryList, nil
}

func (r *Repository) GetLatestUserSyncedRepositories(ctx context.Context, q *db.Queries, user *domain.User) ([]*domain.Repository, error) {
	repositories, err := q.ListStarredRepositoriesBySyncedAt(ctx, db.ListStarredRepositoriesBySyncedAtParams{
		UserID: user.Id,
		Limit:  10,
	})
	if err != nil {
		return nil, err
	}

	repositoryList := make([]*domain.Repository, 0, len(repositories))
	for _, repository := range repositories {

		repo, err := domain.FromRepositoryEntity(
			&domain.RepositoryEntity{
				Id:          repository.ID,
				Name:        repository.Name,
				OwnerName:   repository.OwnerName,
				AvatarURL:   repository.OwnerAvatarUrl,
				HtmlURL:     repository.HtmlUrl,
				Homepage:    *repository.Homepage,
				Description: *repository.Description,
				CreatedAt:   repository.CreatedAt,
				UpdatedAt:   repository.UpdatedAt,
				SyncedAt:    repository.SyncedAt,
				Watchers:    repository.Watchers,
				Forks: 	repository.Forks,
				OpenIssues:  repository.OpenIssues,
				Language:    repository.Language,
				Archived:    repository.Archived,
				Topics:      repository.Topics,
			},
		)

		if err != nil {
			return nil, err
		}

		repositoryList = append(repositoryList, repo)
	}

	return repositoryList, nil
}
