package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type RepoEmbeddingInfo struct {
	ID                   uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RepoID               int64      `gorm:"type:bigint;not null;unique"`
	CreatedAt            time.Time  `gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt            time.Time  `gorm:"type:timestamp with time zone;not null;default:now()"`
	DeletedAt            *time.Time `gorm:"type:timestamp with time zone"`
	FullName             string     `gorm:"type:text;not null"`
	Description          string     `gorm:"type:text"`
	Readme               string     `gorm:"type:text"`
	DescriptionEmbedding *pgvector.Vector  `gorm:"type:vector()"`
	ReadmeEmbedding      *pgvector.Vector  `gorm:"type:vector()"`
}

func (RepoEmbeddingInfo) TableName() string {
	return "repo_embedding_info"
}
