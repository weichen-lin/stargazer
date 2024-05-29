from neo4j import GraphDatabase
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


@dataclass
class RepoVectorizeInfo:
    repo_id: int
    gpt_summary: Optional[str]
    summary_vector: list[float]
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

    def get_user_info(self, email) -> Optional[UserInfo]:
        with self.driver.session() as session:
            record = session.execute_read(self._get_user_info, email)

            if record is None:
                return None

            info = record.data()

            return UserInfo(
                limit=info["limit"],
                openai_key=info["openai_key"],
                cosine=info["cosine"],
            )

    def get_repo_info(self, email: str, repo_id: int) -> Optional[RepoVectorizeInfo]:
        """
        Get the repository information from the graph database
        """
        with self.driver.session() as session:
            record = session.execute_read(self._get_repo_info, email, repo_id)

            if record is None:
                return None

            info = record.data()

            return RepoVectorizeInfo(
                repo_id=info["repo_id"],
                gpt_summary=info["gpt_summary"],
                summary_vector=info["summary_vector"],
                html_url=info["html_url"],
            )

    def save_failed_vectorize_result(self, email: str, repo_id: int, message: str):
        """
        Save the repository information to the graph database
        """
        with self.driver.session() as session:
            session.execute_write(
                self._save_failed_vectorize_result, email, repo_id, message
            )

    def save_repo_info(
        self, email: str, repo_id: int, gpt_summary: str, summary_vector: list[float]
    ):
        """
        Save the repository information to the graph database
        """
        with self.driver.session() as session:
            session.execute_write(
                self._save_repo_info, email, repo_id, gpt_summary, summary_vector
            )

    def get_full_text_repos(self, email: str, query: str):
        with self.driver.session() as session:
            records = session.execute_read(self._get_full_text_repos, email, query)

            return [
                RepoInfo(
                    avatar_url=info["avatar_url"],
                    full_name=info["full_name"],
                    description=info["description"],
                    html_url=info["html_url"],
                    repo_id=info["repo_id"],
                )
                for info in records
            ]

    def get_suggestion_repos(
        self, email: str, limit: int, similarity: float, vector: list[float]
    ):
        with self.driver.session() as session:
            records = session.execute_read(
                self._get_suggestion_repos, email, limit, similarity, vector
            )

            return [
                RepoInfo(
                    avatar_url=info["avatar_url"],
                    full_name=info["full_name"],
                    description=info["description"],
                    html_url=info["html_url"],
                    repo_id=info["repo_id"],
                )
                for info in records
            ]

    @staticmethod
    def _make_index(tx):
        # full-text index for searching repositories
        tx.run(
            """
            CREATE FULLTEXT INDEX REPOSITORY_FULL_TEXT_SEARCH IF NOT EXISTS
            FOR (r:Repository) ON EACH [r.full_name, r.description]
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

    @staticmethod
    def _get_user_info(tx, email: str):
        result = tx.run(
            """
            MATCH (u:User {email: $email})-[:HAS_CONFIG]-(c:Config)
            RETURN c.limit as limit, c.openai_key as openai_key, c.cosine as cosine
            """,
            email=email,
        )

        return result.single()

    @staticmethod
    def _get_repo_info(tx, email: str, repo_id: int):
        result = tx.run(
            """
            MATCH (u:User {email: $email})-[s:STARS]-(r:Repository {repo_id: $repo_id})
            RETURN r.repo_id as repo_id, s.gpt_summary as gpt_summary, s.summary_vector as summary_vector, r.html_url as html_url
            """,
            repo_id=repo_id,
            email=email,
        )

        return result.single()

    @staticmethod
    def _save_repo_info(
        tx, email: str, repo_id: int, gpt_summary: str, summary_vector: list[float]
    ):
        tx.run(
            """
            MATCH (u:User {email: $email})-[s:STARS]-(r:Repository {repo_id: $repo_id})
            SET s += {
                is_vectorized: true,
                gpt_summary: $gpt_summary,
                summary_vector: $summary_vector,
                last_vectorized_at: datetime()
            }
            """,
            email=email,
            repo_id=repo_id,
            gpt_summary=gpt_summary,
            summary_vector=summary_vector,
        )

    @staticmethod
    def _save_failed_vectorize_result(tx, email: str, repo_id: int, reason: str):
        tx.run(
            """
            MATCH (u:User {email: $email})-[s:STARS]-(r:Repository {repo_id: $repo_id})
            SET s += {
                is_vectorized: false,
                reason: $reason,
                last_vectorized_at: datetime()
            }
            """,
            email=email,
            repo_id=repo_id,
            reason=reason,
        )

    @staticmethod
    def _get_suggestion_repos(
        tx, email: str, limit: int, similarity: float, vector: list[float]
    ):
        result = tx.run(
            """
            CALL db.index.vector.queryRelationships("STARS_SUMMARY_VECTOR_INDEX", 5, $vector) YIELD relationship, score
            MATCH (User {email: $email})-[relationship]-(r:Repository)
            WHERE score > $similarity
            RETURN r.repo_id as repo_id, r.avatar_url as avatar_url, r.full_name as full_name, r.html_url as html_url, relationship.gpt_summary as description
            """,
            limit=limit,
            vector=vector,
            email=email,
            similarity=similarity,
        )

        data = list(result.data())

        return data

    @staticmethod
    def _get_full_text_repos(tx, email: str, search_query: str):
        result = tx.run(
            """
            CALL db.index.fulltext.queryNodes("REPOSITORY_FULL_TEXT_SEARCH", $search_query) YIELD node, score
            MATCH (User {email: $email})-[:STARS]-(node)
            RETURN node.repo_id as repo_id ,node.avatar_url as avatar_url, node.full_name as full_name, node.description as description, node.html_url as html_url
            LIMIT 5
            """,
            email=email,
            search_query=search_query,
        )

        data = list(result.data())

        return data
