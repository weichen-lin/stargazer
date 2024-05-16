from pydantic import BaseModel


class RepoEmbeddingInfoSchema(BaseModel):
    repo_id: int
    email: str


class MessageSchema(BaseModel):
    query: str
    email: str
