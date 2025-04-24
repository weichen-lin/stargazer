package domain

import (
	"errors"
	"strings"
	"time"
)

type Repository struct {
	id          int32
	name        string
	ownerName   string
	avatarURL   string
	htmlURL     string
	homepage    string
	description string
	createdAt   time.Time
	updatedAt   time.Time
	watchers    int32
	openIssues  int32
	language    string
	archived    bool
	topics      []string
}

type RepositoryEntity struct {
	Id          int32     `json:"id"`
	Name        string    `json:"name"`
	OwnerName   string    `json:"owner_name"`
	AvatarURL   string    `json:"avatar_url"`
	HtmlURL     string    `json:"html_url"`
	Homepage    string    `json:"homepage"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Watchers    int32     `json:"watchers"`
	OpenIssues  int32     `json:"open_issues"`
	Language    string    `json:"language"`
	Archived    bool      `json:"archived"`
	Topics      []string  `json:"topics"`
}

func (r *Repository) Id() int32 {
	return r.id
}

func (r *Repository) Name() string {
	return r.name
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
	return r.createdAt.UTC()
}

func (r *Repository) UpdatedAt() time.Time {
	return r.updatedAt.UTC()
}

func (r *Repository) Watchers() int32 {
	return r.watchers
}

func (r *Repository) OpenIssues() int32 {
	return r.openIssues
}

func (r *Repository) Language() string {
	return r.language
}

func (r *Repository) Archived() bool {
	return r.archived
}

func (r *Repository) Topics() []string {
	return r.topics
}

func (r *Repository) checkUrl(url string) error {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return errors.New("URL must start with http:// or https://")
	}
	return nil
}

func (r *Repository) setId(id int32) error {
	if id <= 0 {
		return errors.New("repository ID must be positive")
	}
	r.id = id
	return nil
}

func (r *Repository) setName(name string) error {
	if name == "" {
		return errors.New("repository name cannot be empty")
	}
	r.name = name
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

func (r *Repository) setLanguage(lang string) {
	if lang == "" {
		r.language = "Unknown"
	} else {
		r.language = lang
	}
}

func (r *Repository) setTopics(topics []string) {
	r.topics = topics
}

type CreateRepositoryRequest struct {
	Id          int32    `json:"id"`
	Name        string   `json:"name"`
	OwnerName   string   `json:"owner_name"`
	AvatarURL   string   `json:"avatar_url"`
	HtmlURL     string   `json:"html_url"`
	Homepage    string   `json:"homepage"`
	Description string   `json:"description"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	Watchers    int32    `json:"watchers"`
	OpenIssues  int32    `json:"open_issues"`
	Language    string   `json:"language"`
	Archived    bool     `json:"archived"`
	Topics      []string `json:"topics"`
}

func NewRepository(req *CreateRepositoryRequest) (*Repository, error) {
	repo := &Repository{}

	err := repo.setId(req.Id)
	if err != nil {
		return nil, err
	}

	err = repo.setName(req.Name)
	if err != nil {
		return nil, err
	}

	err = repo.setOwnerName(req.OwnerName)
	if err != nil {
		return nil, err
	}

	err = repo.setAvatarURL(req.AvatarURL)
	if err != nil {
		return nil, err
	}

	err = repo.setHTMLURL(req.HtmlURL)
	if err != nil {
		return nil, err
	}

	err = repo.setHomepage(req.Homepage)
	if err != nil {
		return nil, err
	}

	repo.setDescription(req.Description)

	repo.createdAt, _ = time.Parse(time.RFC3339, req.CreatedAt)
	repo.updatedAt, _ = time.Parse(time.RFC3339, req.UpdatedAt)

	repo.watchers = req.Watchers
	repo.openIssues = req.OpenIssues

	repo.setLanguage(req.Language)
	repo.archived = req.Archived
	repo.setTopics(req.Topics)

	return repo, nil
}

func (r *Repository) ToRepositoryEntity() *RepositoryEntity {
	return &RepositoryEntity{
		Id:          r.id,
		Name:        r.name,
		OwnerName:   r.ownerName,
		AvatarURL:   r.avatarURL,
		HtmlURL:     r.htmlURL,
		Homepage:    r.homepage,
		Description: r.description,
		CreatedAt:   r.createdAt.UTC(),
		UpdatedAt:   r.updatedAt.UTC(),
		Watchers:    r.watchers,
		OpenIssues:  r.openIssues,
		Language:    r.language,
		Archived:    r.archived,
		Topics:      r.topics,
	}

}
