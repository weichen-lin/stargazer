from pydantic import BaseModel
from dataclasses import dataclass
from typing import Optional


@dataclass
class UserInfo:
    limit: float
    openai_key: Optional[str]
    cosine: float


@dataclass
class RepoInfo:
    avatar_url: str
    full_name: str
    description: Optional[str]
    html_url: str
    repo_id: int


class RepoEmbeddingInfoSchema(BaseModel):
    repo_id: int
    email: str


class MessageSchema(BaseModel):
    query: str
    email: str
