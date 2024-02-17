package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type RepoEmbeddingInfo struct {
	ID               uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RepoID           int64            `gorm:"type:bigint;not null;unique"`
	CreatedAt        time.Time        `gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt        time.Time        `gorm:"type:timestamp with time zone;not null;default:now()"`
	DeletedAt        *time.Time       `gorm:"type:timestamp with time zone"`
	FullName         string           `gorm:"type:text;not null"`
	AvatarURL        string           `gorm:"type:varying(255)"`
	HTMLURL          string           `gorm:"type:varying(255)"`
	StargazersCount  int              `gorm:"type:integer"`
	Language         string           `gorm:"type:varying(50)"`
	DefaultBranch    string           `gorm:"type:varying(50)"`
	Description      string           `gorm:"type:text"`
	ReadmeSummary    string           `gorm:"type:text"`
	SummaryEmbedding *pgvector.Vector `gorm:"type:vector()"`
}

func (RepoEmbeddingInfo) TableName() string {
	return "repo_embedding_info"
}
