-- CreateExtension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA "public";

-- CreateExtension
CREATE EXTENSION IF NOT EXISTS "vector" WITH SCHEMA "public";

-- CreateTable
CREATE TABLE "repo_embedding_info" (
    "id" UUID NOT NULL DEFAULT uuid_generate_v4(),
    "repo_id" BIGINT NOT NULL,
    "created_at" TIMESTAMPTZ(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMPTZ(6),
    "full_name" TEXT NOT NULL,
    "avatar_url" VARCHAR(255),
    "html_url" VARCHAR(255),
    "stargazers_count" INTEGER,
    "language" VARCHAR(50),
    "default_branch" VARCHAR(50),
    "description" TEXT,
    "readme" TEXT,
    "description_embedding" vector(1536),
    "readme_embedding" vector(1536),

    CONSTRAINT "repo_embedding_info_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "schema_migrations" (
    "version" BIGINT NOT NULL,
    "dirty" BOOLEAN NOT NULL,

    CONSTRAINT "schema_migrations_pkey" PRIMARY KEY ("version")
);

-- CreateIndex
CREATE UNIQUE INDEX "repo_embedding_info_repo_id_key" ON "repo_embedding_info"("repo_id");
