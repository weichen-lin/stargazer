CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "vector";

CREATE TABLE repo_embedding_info (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  repo_id BIGINT UNIQUE,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now(),
  deleted_at TIMESTAMPTZ,
  full_name VARCHAR,
  avatar_url VARCHAR(255),
  html_url VARCHAR(255),
  stargazers_count INT,
  language VARCHAR(50),
  default_branch VARCHAR(50),
  description TEXT,
  readme_summary TEXT,
  summary_embedding VECTOR(150)
);
