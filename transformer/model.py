from pydantic import BaseModel


class RepoEmbeddingInfoSchema(BaseModel):
    repo_id: int
    name: str


class MessageSchema(BaseModel):
    message: str
    name: str
