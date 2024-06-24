import instructor
from openai import OpenAI

from typing import List
from pydantic import Field, BaseModel


class Query(BaseModel):
    """Class representing a single question in a query plan."""

    id: int = Field(..., description="Unique id of the query")
    question: str = Field(
        ...,
        description="Question asked using a question answering system",
    )
    dependencies: List[int] = Field(
        default_factory=list,
        description="List of sub questions that need to be answered before asking this question",
    )


class QueryPlan(BaseModel):
    """Container class representing a tree of questions to ask a question answering system."""

    query_graph: List[Query] = Field(
        ..., description="The query graph representing the plan"
    )

    def _dependencies(self, ids: List[int]) -> List[Query]:
        """Returns the dependencies of a query given their ids."""
        return [q for q in self.query_graph if q.id in ids]


def Planner(question: str, api_key: str) -> QueryPlan:
    PLANNING_MODEL = "gpt-3.5-turbo"

    messages = [
        {
            "role": "system",
            "content": """
                You are a world-class query planning algorithm capable of breaking apart questions into their dependency queries,
                such that the answers can be used to inform the parent question.
                Do not answer the questions.
                Simply provide a correct compute graph with relevant keywords to ask and their dependencies.
                Before you call the function, think step-by-step to get a better understanding of the problem.
            """,
        },
        {
            "role": "user",
            "content": f"Consider: {question}\nGenerate the correct query plan.",
        },
    ]

    client = instructor.from_openai(
        OpenAI(
            api_key=api_key,
        )
    )

    root = client.chat.completions.create(
        model=PLANNING_MODEL,
        temperature=0,
        response_model=QueryPlan,
        messages=messages,
        max_tokens=1000,
    )

    return root
