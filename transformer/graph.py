from neo4j import GraphDatabase
from dataclasses import dataclass
from typing import Optional


@dataclass
class UserInfo:
    limit: int
    openai_key: Optional[str]
    cosine: float


@dataclass
class RepoInfo:
    avatar_url: str
    full_name: str
    description: Optional[str]
    readme_summary: Optional[str]
    html_url: str


@dataclass
class RepoVectorInfo:
    repo_id: int
    repo_vector: list[float]
    readme_summary: Optional[str]
    html_url: str


class Neo4jOperations:
    def __init__(self, uri, user, password):
        self.driver = GraphDatabase.driver(uri, auth=(user, password))

        try:
            with self.driver.session() as session:
                session.execute_write(self._make_index)
        except Exception as e:
            print(f"Error while creating index: {str(e)}")

    def close(self):
        self.driver.close()

    def get_user_info(self, name) -> Optional[UserInfo]:
        with self.driver.session() as session:
            records = session.execute_read(self._get_user_info, name)

            if records is None:
                return None

            info = records[0]

            return UserInfo(
                limit=info["limit"], openai_key=info["openai_key"], cosine=info["cosine"]
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

            return RepoVectorInfo(
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

    def get_suggestion_repos(
        self, name: str, limit: int, similarity: float, vector: list[float]
    ):
        with self.driver.session() as session:
            records = session.execute_read(
                self._get_suggestion_repos, name, limit, similarity, vector
            )

            return [
                RepoInfo(
                    avatar_url=info["avatar_url"],
                    full_name=info["full_name"],
                    description=info["description"],
                    readme_summary=info["readme_summary"],
                    html_url=info["html_url"],
                )
                for info in records
            ]

    def get_full_text_repos(self, name: str, limit: int):
        with self.driver.session() as session:
            records = session.execute_read(self._get_full_text_repos, name, limit)

            return [
                RepoInfo(
                    avatar_url=info["avatar_url"],
                    full_name=info["full_name"],
                    description=info["description"],
                    readme_summary=info["readme_summary"],
                    html_url=info["html_url"],
                )
                for info in records
            ]

    def _make_index(tx):
        # full-text index for searching repositories
        tx.run(
            """
            CREATE FULLTEXT INDEX REPOSITORY_FULL_TEXT_SEARCH IF NOT EXISTS
            FOR (r:Repository) ON EACH [r.full_name, r.description]
        """
        )

        # full-text index for searching summary
        tx.run(
            """
            CREATE FULLTEXT INDEX STARS_SUMMARY_FULL_TEXT_SEARCH IF NOT EXISTS
            FOR ()-[s:STARS]-() ON EACH [s.gpt_summary]
        """
        )

        # vector index for semantic searching repositories
        tx.run(
            """
            CREATE VECTOR INDEX `STARS_SUMMARY_VECTOR_INDEX` IF NOT EXISTS
            FOR ()-[s:STARS]-() ON (s.summary_vector)
            OPTIONS {indexConfig: {
                `vector.dimensions`: 384,
                `vector.similarity_function`: 'cosine'
            }};
        """
        )

    def _get_user_info(tx, email):
        result = tx.run(
            """
            MATCH (u:User {email: $email})-[:HAS_CONFIG]-(c:Config)
            RETURN c
            """, 
            email=email
        )

        return result.single()

    def _get_repo_info(tx, email: str, repo_id: int):
        result = tx.run(
            """
            MATCH (u:User {email: $email})-[s:STARS]-(r:Repository {repo_id: $repo_id})
            RETURN s.gpt_summary as gpt_summary, s.summary_vector as summary_vector, r.html_url as html_url
            """, 
            repo_id=repo_id,
            email=email
            
        )

        return result.single()

    def _save_repo_info(
        tx, repo_id: int, readme_summary: str, repo_vector: list[float]
    ):
        tx.run(
            "MATCH (r:Repository { repo_id: $repo_id }) SET r.readme_summary = $readme_summary, r.repo_vector = $repo_vector",
            repo_id=repo_id,
            readme_summary=readme_summary,
            repo_vector=repo_vector,
        )

    def _get_suggestion_repos(
        tx, name: str, limit: int, similarity: float, vector: list[float]
    ):
        result = tx.run(
            """
            CALL db.index.vector.queryNodes("REPOSITORY_VECTOR_INDEX", 5, $vector)
            YIELD node, score
            MATCH (User {name: $name})-[:STARS]-(node)
            WHERE score > $similarity
            RETURN 
            node.avatar_url as avatar_url,
            node.full_name as full_name, 
            node.description as description, 
            node.readme_summary as readme_summary, 
            node.html_url as html_url
            """,
            limit=limit,
            vector=vector,
            name=name,
            similarity=similarity,
        )

        data = list(result.data())

        return data

    def _get_full_text_repos(tx, message: str, name: str):
        result = tx.run(
            """
            CALL db.index.fulltext.queryNodes("REPOSITORY_FULL_TEXT_SEARCH", $message) YIELD node, score
            MATCH (User {name: $name})-[:STARS]-(node)
            RETURN node.avatar_url as avatar_url, node.full_name as full_name, node.description as description, node.readme_summary as readme_summary
            LIMIT 5
            """,
            massage=message,
            name=name,
        )

        data = list(result.data())

        return data
