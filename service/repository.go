package service

import (
	"context"

	"github.com/weichen-lin/stargazer/db"
	"github.com/weichen-lin/stargazer/domain"
	"github.com/weichen-lin/stargazer/repository"
)

func (s *Service) GetLanguageDistribution(ctx context.Context) ([]*repository.LanguageDistribution, error) {
	languageDistributionList := make([]*repository.LanguageDistribution, 0)

	err := s.db.Run(ctx, func(q *db.Queries) error {
		user, err := s.repository.GetUserByClerkId(ctx, q)
		if err != nil {
			return err
		}

		languageDistributionList, err = s.repository.GetLanguageDistribution(ctx, q, user)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return languageDistributionList, nil
}

func (s *Service) GetLatestUserUpdatedRepositories(ctx context.Context) ([]*domain.Repository, error) {
	repos := make([]*domain.Repository, 0)

	err := s.db.Run(ctx, func(q *db.Queries) error {
		user, err := s.repository.GetUserByClerkId(ctx, q)
		if err != nil {
			return err
		}

		repos, err = s.repository.GetLatestUserUpdatedRepositories(ctx, q, user)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return repos, nil
}

func (s *Service) GetLatestUserSyncedRepositories(ctx context.Context) ([]*domain.Repository, error) {
	repos := make([]*domain.Repository, 0)

	err := s.db.Run(ctx, func(q *db.Queries) error {
		user, err := s.repository.GetUserByClerkId(ctx, q)
		if err != nil {
			return err
		}

		repos, err = s.repository.GetLatestUserSyncedRepositories(ctx, q, user)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return repos, nil
}
