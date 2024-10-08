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
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	FullName        string   `json:"full_name"`
	Owner           Owner    `json:"owner"`
	HTMLURL         string   `json:"html_url"`
	Description     string   `json:"description"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
	Homepage        string   `json:"homepage"`
	StargazersCount int64    `json:"stargazers_count"`
	WatchersCount   int64    `json:"watchers_count"`
	Language        string   `json:"language"`
	Archived        bool     `json:"archived"`
	OpenIssuesCount int64    `json:"open_issues_count"`
	Topics          []string `json:"topics"`
	DefaultBranch   string   `json:"default_branch"`
}

type RepositoryEntity struct {
	RepoID            int64    `json:"repo_id"`
	RepoName          string   `json:"repo_name"`
	OwnerName         string   `json:"owner_name"`
	AvatarURL         string   `json:"avatar_url"`
	HtmlURL           string   `json:"html_url"`
	Homepage          string   `json:"homepage"`
	Description       string   `json:"description"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
	StargazersCount   int64    `json:"stargazers_count"`
	Language          string   `json:"language"`
	WatchersCount     int64    `json:"watchers_count"`
	OpenIssuesCount   int64    `json:"open_issues_count"`
	DefaultBranch     string   `json:"default_branch"`
	Archived          bool     `json:"archived"`
	Topics            []string `json:"topics"`
	ExternalCreatedAt string   `json:"external_created_at"`
	LastSyncedAt      string   `json:"last_synced_at"`
	LastModifiedAt    string   `json:"last_modified_at"`
}

type Repository struct {
	repoID            int64
	repoName          string
	ownerName         string
	avatarURL         string
	htmlURL           string
	homepage          string
	description       string
	createdAt         time.Time
	updatedAt         time.Time
	stargazersCount   int64
	watchersCount     int64
	openIssuesCount   int64
	language          string
	defaultBranch     string
	archived          bool
	topics            []string
	externalCreatedAt time.Time
	lastSyncedAt      time.Time
	lastModifiedAt    time.Time
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

func (r *Repository) StargazersCount() int64 {
	return r.stargazersCount
}

func (r *Repository) WatchersCount() int64 {
	return r.watchersCount
}

func (r *Repository) OpenIssuesCount() int64 {
	return r.openIssuesCount
}

func (r *Repository) Language() string {
	return r.language
}

func (r *Repository) DefaultBranch() string {
	return r.defaultBranch
}

func (r *Repository) Archived() bool {
	return r.archived
}

func (r *Repository) Topics() []string {
	return r.topics
}

func (r *Repository) ExternalCreateAt() time.Time {
	return r.externalCreatedAt
}

func (r *Repository) LastSyncedAt() time.Time {
	return r.lastSyncedAt
}

func (r *Repository) LastModifiedAt() time.Time {
	return r.lastModifiedAt
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
	if url == "" {
		r.homepage = url
		return nil
	}

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

func (r *Repository) setCreatedAt(t string) error {
	if t == "" {
		return errors.New("created time cannot be empty")
	}

	parsedTime, err := ParseTime(t)
	if err != nil {
		return err
	}

	r.createdAt = parsedTime
	return nil
}

func (r *Repository) setUpdatedAt(t string) error {
	if t == "" {
		return errors.New("updated time cannot be empty")
	}

	parsedTime, err := ParseTime(t)
	if err != nil {
		return err
	}

	if parsedTime.Before(r.createdAt) {
		return errors.New("updated time cannot be before created time")
	}

	r.updatedAt = parsedTime
	return nil
}

func (r *Repository) setStargazersCount(count int64) error {
	if count < 0 {
		return errors.New("stargazers count cannot be negative")
	}
	r.stargazersCount = count
	return nil
}

func (r *Repository) setWatchersCount(count int64) error {
	if count < 0 {
		return errors.New("watchers count cannot be negative")
	}
	r.watchersCount = count
	return nil
}

func (r *Repository) setOpenIssuesCount(count int64) error {
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

func (r *Repository) setTopics(topics []string) {
	r.topics = topics
}

func (r *Repository) setExternalCreatedAt(t time.Time) error {
	if t.IsZero() {
		return errors.New("external created time cannot be empty")
	}

	r.externalCreatedAt = t
	return nil
}

func (r *Repository) setLastSyncedAt(t time.Time) error {
	if t.IsZero() {
		return errors.New("last sync time cannot be empty")
	}

	r.lastSyncedAt = t
	return nil
}

func (r *Repository) setModifiedAt(t time.Time) error {
	if t.IsZero() {
		return errors.New("last modified time cannot be empty")
	}

	r.lastModifiedAt = t
	return nil
}

func (r *Repository) ToRepositoryEntity() *RepositoryEntity {
	layout := "2006-01-02T15:04:05Z07:00"

	return &RepositoryEntity{
		RepoID:            r.repoID,
		RepoName:          r.repoName,
		OwnerName:         r.ownerName,
		AvatarURL:         r.avatarURL,
		HtmlURL:           r.htmlURL,
		Homepage:          r.homepage,
		Description:       r.description,
		CreatedAt:         r.createdAt.Format(layout),
		UpdatedAt:         r.updatedAt.Format(layout),
		StargazersCount:   r.stargazersCount,
		WatchersCount:     r.watchersCount,
		OpenIssuesCount:   r.openIssuesCount,
		Language:          r.language,
		DefaultBranch:     r.defaultBranch,
		Archived:          r.archived,
		Topics:            r.topics,
		ExternalCreatedAt: r.externalCreatedAt.Format(layout),
		LastSyncedAt:      r.lastSyncedAt.Format(layout),
		LastModifiedAt:    r.lastModifiedAt.Format(layout),
	}
}

func FromRepositoryEntity(repositoryEntity *RepositoryEntity) (*Repository, error) {
	r := &Repository{}

	if err := r.setRepoID(int64(repositoryEntity.RepoID)); err != nil {
		return nil, err
	}

	if err := r.setRepoName(repositoryEntity.RepoName); err != nil {
		return nil, err
	}

	if err := r.setOwnerName(repositoryEntity.OwnerName); err != nil {
		return nil, err
	}

	if err := r.setAvatarURL(repositoryEntity.AvatarURL); err != nil {
		return nil, err
	}

	if err := r.setHTMLURL(repositoryEntity.HtmlURL); err != nil {
		return nil, err
	}

	if err := r.setHomepage(repositoryEntity.Homepage); err != nil {
		return nil, err
	}

	r.setDescription(repositoryEntity.Description)

	if err := r.setCreatedAt(repositoryEntity.CreatedAt); err != nil {
		return nil, err
	}

	if err := r.setUpdatedAt(repositoryEntity.UpdatedAt); err != nil {
		return nil, err
	}

	if err := r.setStargazersCount(repositoryEntity.StargazersCount); err != nil {
		return nil, err
	}

	if err := r.setWatchersCount(repositoryEntity.WatchersCount); err != nil {
		return nil, err
	}

	if err := r.setOpenIssuesCount(repositoryEntity.OpenIssuesCount); err != nil {
		return nil, err
	}

	r.setLanguage(repositoryEntity.Language)

	if err := r.setDefaultBranch(repositoryEntity.DefaultBranch); err != nil {
		return nil, err
	}

	r.setArchived(repositoryEntity.Archived)
	r.setTopics(repositoryEntity.Topics)

	externalTime, err := ParseTime(repositoryEntity.ExternalCreatedAt)
	if err != nil {
		return nil, err
	}

	r.setExternalCreatedAt(externalTime)

	lastSyncedAtTime, err := ParseTime(repositoryEntity.LastSyncedAt)
	if err != nil {
		return nil, err
	}

	r.setLastSyncedAt(lastSyncedAtTime)

	lastModifiedAtTime, err := ParseTime(repositoryEntity.LastModifiedAt)
	if err != nil {
		return nil, err
	}

	r.setModifiedAt(lastModifiedAtTime)

	return r, nil
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

	repo.setTopics(githubRepo.Topics)

	repo.setExternalCreatedAt(time.Now())
	repo.setLastSyncedAt(time.Now())
	repo.setModifiedAt(time.Now())

	return repo, nil
}
