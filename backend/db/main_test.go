package db

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/stargazer/domain"
	"golang.org/x/exp/rand"
)

var db *Database

func NewTestDatabase() *Database {
	driver, err := neo4j.NewDriverWithContext(
		"neo4j://localhost:7687",
		neo4j.BasicAuth("neo4j", "password", ""),
	)

	if err != nil {
		panic(err)
	}

	return &Database{
		Driver: driver,
	}
}

func createFakeUser(t *testing.T) (*domain.User, context.Context) {
	entity := &domain.UserEntity{
		Name:              faker.Name(),
		Email:             faker.Email(),
		Image:             faker.URL(),
		AccessToken:       faker.Sentence(),
		Provider:          faker.Sentence(),
		ProviderAccountId: faker.Sentence(),
		Scope:             faker.Sentence(),
		AuthType:          faker.Sentence(),
		TokenType:         faker.Sentence(),
	}

	user := domain.FromUserEntity(entity)

	err := db.CreateUser(user)
	require.NoError(t, err)

	ctx, err := WithEmail(context.Background(), user.Email())
	require.NoError(t, err)
	require.Equal(t, ctx.Value("email"), user.Email())

	return user, ctx
}

var Languages []string = []string{
	"Go",
	"Python",
	"JavaScript",
	"Java",
	"C++",
	"C#",
	"Ruby",
	"Swift",
	"Kotlin",
	"TypeScript",
	"Rust",
	"PHP",
	"Scala",
	"Haskell",
	"Perl",
	"R",
	"MATLAB",
	"Dart",
	"Lua",
	"Groovy",
}

func getRandomLanguage() string {
	seed := uint64(time.Now().UnixNano()) & ((1 << 63) - 1)
	rand.Seed(seed)

	randomIndex := rand.Intn(len(Languages))

	return Languages[randomIndex]
}

func createRepositoryAtFakeUser(t *testing.T, user *domain.User) *domain.Repository {
	repositoryEntity := &domain.RepositoryEntity{
		RepoID:            faker.RandomUnixTime(),
		RepoName:          faker.Name(),
		OwnerName:         faker.Name(),
		AvatarURL:         faker.URL(),
		HtmlURL:           faker.URL(),
		Homepage:          faker.URL(),
		CreatedAt:         "2024-01-01T00:00:00Z",
		UpdatedAt:         "2024-01-02T00:00:00Z",
		StargazersCount:   100,
		WatchersCount:     50,
		OpenIssuesCount:   10,
		DefaultBranch:     "main",
		Description:       faker.Sentence(),
		Language:          getRandomLanguage(),
		Archived:          false,
		ExternalCreatedAt: time.Now().Format(time.RFC3339),
		LastSyncedAt:      time.Now().Format(time.RFC3339),
		LastModifiedAt:    time.Now().Format(time.RFC3339),
	}

	repo, err := domain.FromRepositoryEntity(repositoryEntity)
	require.NoError(t, err)
	require.NotEmpty(t, repo)

	ctx, err := WithEmail(context.Background(), user.Email())
	require.NoError(t, err)
	require.NotEmpty(t, ctx)

	err = db.CreateRepository(ctx, repo)
	require.NoError(t, err)

	return repo
}

func TestMain(m *testing.M) {
	db = NewTestDatabase()

	os.Exit(m.Run())
}
