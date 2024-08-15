package domain

import (
	"errors"
	"strings"
	"time"
)

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

type Repository struct {
	repoID          int64
	repoName        string
	ownerName       string
	avatarURL       string
	htmlURL         string
	homepage        string
	description     string
	createdAt       time.Time
	updatedAt       time.Time
	stargazersCount int
	watchersCount   int
	openIssuesCount int
	language        string
	defaultBranch   string
	archived        bool
}

func (r *Repository) RepoID() int64 {
	return r.repoID
}

func (r *Repository) RepoName() string {
	return r.repoName
}

func (r *Repository) OwnerName() string {
	return r.ownerName
}

func (r *Repository) AvatarURL() string {
	return r.avatarURL
}

func (r *Repository) HTMLURL() string {
	return r.htmlURL
}

func (r *Repository) Homepage() string {
	return r.homepage
}

func (r *Repository) Description() string {
	return r.description
}

func (r *Repository) CreatedAt() time.Time {
	return r.createdAt
}

func (r *Repository) UpdatedAt() time.Time {
	return r.updatedAt
}

func (r *Repository) StargazersCount() int {
	return r.stargazersCount
}

func (r *Repository) WatchersCount() int {
	return r.watchersCount
}

func (r *Repository) OpenIssuesCount() int {
	return r.openIssuesCount
}

func (r *Repository) Language() string {
	return r.language
}

func (r *Repository) DefaultBranch() string {
	return r.defaultBranch
}

func (r *Repository) IsArchived() bool {
	return r.archived
}

func (r *Repository) checkUrl(url string) error {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return errors.New("URL must start with http:// or https://")
	}
	return nil
}

func (r *Repository) setRepoID(id int64) error {
	if id <= 0 {
		return errors.New("repository ID must be positive")
	}
	r.repoID = id
	return nil
}

func (r *Repository) setRepoName(name string) error {
	if name == "" {
		return errors.New("repository name cannot be empty")
	}
	r.repoName = name
	return nil
}

func (r *Repository) setOwnerName(name string) error {
	if name == "" {
		return errors.New("owner name cannot be empty")
	}
	r.ownerName = name
	return nil
}

func (r *Repository) setAvatarURL(url string) error {
	err := r.checkUrl(url)
	if err != nil {
		return err
	}
	r.avatarURL = url
	return nil
}

func (r *Repository) setHTMLURL(url string) error {
	err := r.checkUrl(url)
	if err != nil {
		return err
	}
	r.htmlURL = url
	return nil
}

func (r *Repository) setHomepage(url string) error {
	err := r.checkUrl(url)
	if err != nil {
		return err
	}
	r.homepage = url
	return nil
}

func (r *Repository) setDescription(desc string) {
	r.description = desc
}

func (r *Repository) setCreatedAt(t time.Time) error {
	if t.IsZero() {
		return errors.New("created time cannot be empty")
	}
	r.createdAt = t
	return nil
}

func (r *Repository) setUpdatedAt(t time.Time) error {
	if t.IsZero() {
		return errors.New("updated time cannot be empty")
	}

	if t.Before(r.createdAt) {
		return errors.New("updated time cannot be before created time")
	}

	r.updatedAt = t
	return nil
}

func (r *Repository) setStargazersCount(count int) error {
	if count < 0 {
		return errors.New("stargazers count cannot be negative")
	}
	r.stargazersCount = count
	return nil
}

func (r *Repository) setWatchersCount(count int) error {
	if count < 0 {
		return errors.New("watchers count cannot be negative")
	}
	r.watchersCount = count
	return nil
}

func (r *Repository) setOpenIssuesCount(count int) error {
	if count < 0 {
		return errors.New("open issues count cannot be negative")
	}
	r.openIssuesCount = count
	return nil
}

func (r *Repository) setLanguage(lang string) {
	if lang == "" {
		r.language = "Unknown"
	} else {
		r.language = lang
	}
}

func (r *Repository) setDefaultBranch(branch string) error {
	if branch == "" {
		return errors.New("default branch cannot be empty")
	}
	r.defaultBranch = branch
	return nil
}

func (r *Repository) setArchived(archived bool) {
	r.archived = archived
}

func NewRepository(githubRepo *GithubRepository) (*Repository, error) {
	repo := &Repository{}

	if err := repo.setRepoID(int64(githubRepo.ID)); err != nil {
		return nil, err
	}

	if err := repo.setRepoName(githubRepo.Name); err != nil {
		return nil, err
	}

	if err := repo.setOwnerName(githubRepo.Owner.Login); err != nil {
		return nil, err
	}

	if err := repo.setAvatarURL(githubRepo.Owner.AvatarURL); err != nil {
		return nil, err
	}

	if err := repo.setHTMLURL(githubRepo.HTMLURL); err != nil {
		return nil, err
	}

	if err := repo.setHomepage(githubRepo.Homepage); err != nil {
		return nil, err
	}

	repo.setDescription(githubRepo.Description)

	if err := repo.setCreatedAt(githubRepo.CreatedAt); err != nil {
		return nil, err
	}

	if err := repo.setUpdatedAt(githubRepo.UpdatedAt); err != nil {
		return nil, err
	}

	if err := repo.setStargazersCount(githubRepo.StargazersCount); err != nil {
		return nil, err
	}

	if err := repo.setWatchersCount(githubRepo.WatchersCount); err != nil {
		return nil, err
	}

	if err := repo.setOpenIssuesCount(githubRepo.OpenIssuesCount); err != nil {
		return nil, err
	}

	repo.setLanguage(githubRepo.Language)

	if err := repo.setDefaultBranch(githubRepo.DefaultBranch); err != nil {
		return nil, err
	}

	repo.setArchived(githubRepo.Archived)

	return repo, nil
}
