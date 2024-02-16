from flask_sqlalchemy import SQLAlchemy
from sqlalchemy.dialects.postgresql import UUID, ARRAY
from sqlalchemy import func
from pydantic import BaseModel
from pgvector.sqlalchemy import Vector

db = SQLAlchemy()

class RepoEmbeddingInfo(db.Model):
    __tablename__ = 'repo_embedding_info'

    id = db.Column(UUID(as_uuid=True), primary_key=True, server_default=func.uuid_generate_v4())
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

    def __repr__(self):
        return '<RepoEmbeddingInfo %r>' % self.repo_id


class RepoEmbeddingInfoSchema(BaseModel):
    repo_id: int