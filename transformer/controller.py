from playwright.sync_api import sync_playwright
from config import db
from helper import get_token_length, get_embedding, get_formatted_text
from openai import OpenAI, AuthenticationError, APIConnectionError
from plan import Planner
from cleaner import flatten_and_deduplicate


def Crawler(id: int, email: str) -> tuple[str, int]:

    info = db.get_user_info(email)
    if info is None:
        raise ValueError(f"{email} not found")

    if info.openai_key is None or info.openai_key == "":
        raise ValueError(f"{email} does not have an valid openAIKey")

    client = OpenAI(
        api_key=info.openai_key,
    )

    repo = db.get_repo_info(email, id)

    if repo is None:
        raise ValueError(f"Repo with id {id} not found")

    if repo.html_url is None:
        raise ValueError(f"Repo with id {id} does not have an html_url")

    if repo.summary_vector is not None:
        return f"already generate embedding on repo: {id}", 201

    if repo.gpt_summary is None or repo.gpt_summary == "":
        with sync_playwright() as p:
            try:
                browser = p.chromium.launch(headless=True)
                page = browser.new_page()
                page.goto(repo.html_url)

                selectors = ["article.markdown-body", ".plain"]

                article = page.wait_for_selector(f'{", ".join(selectors)}')
                content = f"The following is a repository from github. After reading the material, can you give me a summary? Reply in English\n{article.text_content()}"

                tokens = get_token_length(content)

                chat_completion = client.chat.completions.create(
                    messages=[
                        {
                            "role": "user",
                            "content": content,
                        }
                    ],
                    model="gpt-3.5-turbo" if tokens < 16000 else "gpt-4-turbo-preview",
                )

                summary = chat_completion.choices[0].message.content

                repo.gpt_summary = summary or ""

                summary = summary.replace("\n", " ")
                vector = get_embedding(summary)
                repo.summary_vector = vector

                try:
                    db.save_repo_info(
                        email=email,
                        repo_id=id,
                        gpt_summary=repo.gpt_summary,
                        summary_vector=repo.summary_vector,
                    )

                    return f"success generate embedding on repo: {id}", 200

                except Exception as e:
                    raise ValueError(
                        f"Error while inserting to neo4j: {str(id)}, {str(e)}"
                    )

            except AuthenticationError:
                raise ValueError("Invalid OpenAI Key")

            except APIConnectionError:
                raise ValueError("API Connection Error")

            except Exception as e:
                db.save_failed_vectorize_result(email, id, str(e))
                raise ValueError(f"Error while crawling: {str(e)}")


def VectorSearcher(email: str, query: str) -> list[dict]:
    info = db.get_user_info(email)
    if info is None:
        raise ValueError(f"User {email} not found")

    if info.openai_key is None or info.openai_key == "":
        raise ValueError(f"User {email} does not have a valid openAIKey")

    client = OpenAI(api_key=info.openai_key)
    client.models.list()

    plans = Planner(query, info.openai_key)

    vectors = [get_embedding(get_formatted_text(q.question)) for q in plans.query_graph]
    repos = [
        db.get_suggestion_repos(email, info.limit, info.cosine, vector)
        for vector in vectors
    ]

    return flatten_and_deduplicate(repos), 200


def FullTextSearcher(email: str, query: str) -> list[dict]:

    info = db.get_user_info(email)
    if info is None:
        raise ValueError(f"User {email} not found")

    repos = db.get_full_text_repos(email, query)

    return repos, 200
