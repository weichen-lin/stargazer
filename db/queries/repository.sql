-- name: UpsertOwner :one
INSERT INTO owner (id, name, avatar_url)
VALUES ($1, $2, $3)
ON CONFLICT (id) DO UPDATE
SET name = EXCLUDED.name,
    avatar_url = EXCLUDED.avatar_url
RETURNING *;

-- name: UpsertRepository :one
INSERT INTO repositories (id, name, owner_id, html_url, homepage, description, watchers, forks, open_issues, language, archived, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
ON CONFLICT (id) DO UPDATE
SET name = EXCLUDED.name,
    owner_id = EXCLUDED.owner_id,
    html_url = EXCLUDED.html_url,
    homepage = EXCLUDED.homepage,
    description = EXCLUDED.description,
    watchers = EXCLUDED.watchers,
    forks = EXCLUDED.forks,
    open_issues = EXCLUDED.open_issues,
    language = EXCLUDED.language,
    archived = EXCLUDED.archived,
    created_at = EXCLUDED.created_at,
    updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: UpsertTopics :one
INSERT INTO topics (name)
VALUES ($1)
ON CONFLICT (name) DO UPDATE
SET name = EXCLUDED.name
RETURNING *;

-- name: UpsertRepositoryTopics :many
INSERT INTO repository_topics (repo_id, topic_id)
VALUES ($1, $2)
ON CONFLICT (repo_id, topic_id) DO UPDATE
SET repo_id = EXCLUDED.repo_id,
    topic_id = EXCLUDED.topic_id
RETURNING *;