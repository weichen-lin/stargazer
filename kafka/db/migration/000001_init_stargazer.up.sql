CREATE EXTENSION IF NOT EXISTS "vector";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS repo_embedding_info (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    repo_id BIGINT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    full_name TEXT NOT NULL,
    avatar_url VARCHAR(255),
    html_url VARCHAR(255),
    stargazers_count INT,
    language VARCHAR(50),
    default_branch VARCHAR(50),
    description TEXT,
    readme TEXT,
    description_embedding VECTOR(1536),
    readme_embedding VECTOR(1536)
);