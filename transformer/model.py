from flask_sqlalchemy import SQLAlchemy
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy import func
from pydantic import BaseModel
from pgvector.sqlalchemy import Vector

db = SQLAlchemy()

class ElasticSearchDoc(BaseModel):
    full_name: str
    avatar_url: str
    html_url: str
    stargazers_count: int
    language: str
    description: str
    readme_summary: str
    elk_vector: list

class RepoEmbeddingInfo(db.Model):
    __tablename__ = "repo_embedding_info"

    id = db.Column(
        UUID(as_uuid=True), primary_key=True, server_default=func.uuid_generate_v4()
    )
    repo_id = db.Column(db.BigInteger, unique=True)
    created_at = db.Column(db.TIMESTAMP(timezone=True), server_default=func.now())
    updated_at = db.Column(db.TIMESTAMP(timezone=True), server_default=func.now())
    deleted_at = db.Column(db.TIMESTAMP(timezone=True))
    full_name = db.Column(db.String)
    avatar_url = db.Column(db.String(255))
    html_url = db.Column(db.String(255))
    stargazers_count = db.Column(db.Integer)
    language = db.Column(db.String(50))
    default_branch = db.Column(db.String(50))
    description = db.Column(db.Text)
    readme_summary = db.Column(db.Text)
    summary_embedding = db.Column(Vector, nullable=True)
    elk_vector = db.Column(Vector, nullable=True)

    def __repr__(self):
        return "<RepoEmbeddingInfo %r>" % self.repo_id

    def _to_json(self):
        return {
            "id": self.id,
            "repo_id": self.repo_id,
            "full_name": self.full_name,
            "avatar_url": self.avatar_url,
            "html_url": self.html_url,
            "stargazers_count": self.stargazers_count,
            "default_branch": self.default_branch,
            "description": self.description,
        }

    def _to_elastic(self) -> ElasticSearchDoc:
        return {
            "full_name": self.full_name,
            "avatar_url": self.avatar_url,
            "html_url": self.html_url,
            "stargazers_count": self.stargazers_count,
            "language": self.language,
            "description": self.description,
            "readme_summary": self.readme_summary,
            "elk_vector": self.elk_vector.tolist() if self.elk_vector is not None else None,
        }

class RepoEmbeddingInfoSchema(BaseModel):
    repo_id: int
    name: str


class MessageSchema(BaseModel):
    message: str
    name: str
