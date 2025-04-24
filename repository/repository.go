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
