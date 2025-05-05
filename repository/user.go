package repository

import (
	"context"
	"net/http"
	"time"

	"github.com/weichen-lin/stargazer/db"
	"github.com/weichen-lin/stargazer/domain"
	"github.com/weichen-lin/stargazer/infrastructure"
)

func (r *Repository) GetUserByClerkId(ctx context.Context, q *db.Queries) (*domain.User, error) {
	clerkId, ok := infrastructure.GetClerkId(ctx)
	if !ok {
		return nil, infrastructure.NewServiceError(
			ctx,
			http.StatusInternalServerError,
			"not found clerk id at context",
			"not found clerk id at context",
		)
	}

	dbUser, err := q.GetUserByClerkId(ctx, clerkId)
	if err != nil {
		return nil, err
	}

	user, err := domain.NewUser(
		dbUser.ID.String(),
		dbUser.ClerkID,
		dbUser.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) SaveUserStarredRepository(ctx context.Context, q *db.Queries, repository_id int32) error {
	clerkId, ok := infrastructure.GetClerkId(ctx)
	if !ok {
		return infrastructure.NewServiceError(
			ctx,
			http.StatusInternalServerError,
			"not found clerk id at context",
			"not found clerk id at context",
		)
	}

	dbUser, err := q.GetUserByClerkId(ctx, clerkId)
	if err != nil {
		return err
	}

	_, err = q.UpsertUserStarRepository(ctx, db.UpsertUserStarRepositoryParams{
		UserID:    dbUser.ID,
		RepoID:    repository_id,
		IsDelete:  false,
		CreatedAt: time.Now().UTC(),
		SyncedAt:  time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		return infrastructure.NewServiceError(
			ctx,
			http.StatusInternalServerError,
			"failed to upsert user star repository",
			err.Error(),
		)
	}

	return nil
}

func (r *Repository) SaveUserCrontab(ctx context.Context, q *db.Queries, repository_counts int32) error {
	clerkId, ok := infrastructure.GetClerkId(ctx)
	if !ok {
		return infrastructure.NewServiceError(
			ctx,
			http.StatusInternalServerError,
			"not found clerk id at context",
			"not found clerk id at context",
		)
	}

	dbUser, err := q.GetUserByClerkId(ctx, clerkId)
	if err != nil {
		return err
	}

	_, err = q.UpsertCrontab(ctx, db.UpsertCrontabParams{
		UserID:     dbUser.ID,
		Stargazers: repository_counts,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	})

	if err != nil {
		return infrastructure.NewServiceError(
			ctx,
			http.StatusInternalServerError,
			"failed to upsert user crontab",
			err.Error(),
		)
	}

	return nil
}

func (r *Repository) GetUserCrontab(ctx context.Context, q *db.Queries) (*domain.Crontab, error) {
	clerkId, ok := infrastructure.GetClerkId(ctx)
	if !ok {
		return nil, infrastructure.NewServiceError(
			ctx,
			http.StatusInternalServerError,
			"not found clerk id at context",
			"not found clerk id at context",
		)
	}

	dbUser, err := q.GetUserByClerkId(ctx, clerkId)
	if err != nil {
		return nil, err
	}

	dbCrontab, err := q.GetCrontab(ctx, dbUser.ID)

	crontab := domain.FromCrontabEntity(&domain.CrontabEntity{
		UserId:     dbCrontab.UserID,
		Stargazers: dbCrontab.Stargazers,
		CreatedAt:  dbCrontab.CreatedAt,
		UpdatedAt:  dbCrontab.UpdatedAt,
	})
	if err != nil {
		return nil, infrastructure.NewServiceError(
			ctx,
			http.StatusInternalServerError,
			"failed to convert crontab entity",
			err.Error(),
		)
	}

	return crontab, nil
}
