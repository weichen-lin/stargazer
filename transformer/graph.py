from neo4j import GraphDatabase
from dataclasses import dataclass
from typing import Optional


@dataclass
class UserInfo:
    limit: int
    openAIKey: Optional[str]
    cosine: float


@dataclass
class RepoInfo:
    repo_id: int
    readme_summary: Optional[str]
    repo_vector: Optional[list[float]]
    html_url: str


class Neo4jOperations:
    def __init__(self, uri, user, password):
        self.driver = GraphDatabase.driver(uri, auth=(user, password))

    def close(self):
        self.driver.close()

    def get_user_info(self, name) -> Optional[UserInfo]:
        with self.driver.session() as session:
            records = session.execute_read(self._get_user_info, name)

            if records is None:
                return None

            info = records[0]

            return UserInfo(
                limit=info["limit"], openAIKey=info["openAIKey"], cosine=info["cosine"]
            )

    def get_repo_info(self, repo_id: int) -> Optional[RepoInfo]:
        """
        Get the repository information from the graph database
        """
        with self.driver.session() as session:
            records = session.execute_read(self._get_repo_info, repo_id)

            if records is None:
                return None

            info = records[0]

            return RepoInfo(
                repo_id=info["repo_id"],
                readme_summary=info["readme_summary"],
                repo_vector=info["repo_vector"],
                html_url=info["html_url"],
            )

    def save_repo_info(
        self, repo_id: int, readme_summary: str, repo_vector: list[float]
    ):
        """
        Save the repository information to the graph database
        """
        with self.driver.session() as session:
            session.execute_write(
                self._save_repo_info, repo_id, readme_summary, repo_vector
            )

    @staticmethod
    def _get_user_info(tx, name):
        result = tx.run("MATCH (u:User { name: $name }) RETURN u", name=name)

        return result.single()

    @staticmethod
    def _get_repo_info(tx, repo_id: int):
        result = tx.run(
            "MATCH (r:Repository { repo_id: $repo_id }) RETURN r", repo_id=repo_id
        )

        return result.single()

    @staticmethod
    def _save_repo_info(
        tx, repo_id: int, readme_summary: str, repo_vector: list[float]
    ):
        tx.run(
            "MATCH (r:Repository { repo_id: $repo_id }) SET r.readme_summary = $readme_summary, r.repo_vector = $repo_vector",
            repo_id=repo_id,
            readme_summary=readme_summary,
            repo_vector=repo_vector,
        )
