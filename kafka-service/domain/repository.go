package domain

import "time"

type Repository struct {
	RepoID          int64     `json:"id"`
	RepoName        string    `json:"repo_name"`
	OwnerName       string    `json:"owner_name"`
	AvatarURL       string    `json:"avatar_url"`
	HTMLURL         string    `json:"html_url"`
	Homepage        string    `json:"homepage"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	StargazersCount int       `json:"stargazers_count"`
	WatchersCount   int       `json:"watchers_count"`
	OpenIssuesCount int       `json:"open_issues_count"`
	Language        string    `json:"language"`
	DefaultBranch   string    `json:"default_branch"`
	Archived        bool      `json:"archived"`
}

type Owner struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
}

type GithubRepository struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Owner           Owner     `json:"owner"`
	HTMLURL         string    `json:"html_url"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Homepage        string    `json:"homepage"`
	StargazersCount int       `json:"stargazers_count"`
	WatchersCount   int       `json:"watchers_count"`
	Language        string    `json:"language"`
	Archived        bool      `json:"archived"`
	OpenIssuesCount int       `json:"open_issues_count"`
	Topics          []string  `json:"topics"`
	DefaultBranch   string    `json:"default_branch"`
}

func NewRepository(githubRepo *GithubRepository) *Repository {
	return &Repository{
		RepoID:          int64(githubRepo.ID),
		RepoName:        githubRepo.Name,
		OwnerName:       githubRepo.Owner.Login,
		AvatarURL:       githubRepo.Owner.AvatarURL,
		HTMLURL:         githubRepo.HTMLURL,
		Homepage:        githubRepo.Homepage,
		Description:     githubRepo.Description,
		CreatedAt:       githubRepo.CreatedAt,
		UpdatedAt:       githubRepo.UpdatedAt,
		StargazersCount: githubRepo.StargazersCount,
		WatchersCount:   githubRepo.WatchersCount,
		OpenIssuesCount: githubRepo.OpenIssuesCount,
		Language:        githubRepo.Language,
		DefaultBranch:   githubRepo.DefaultBranch,
		Archived:        githubRepo.Archived,
	}
}
