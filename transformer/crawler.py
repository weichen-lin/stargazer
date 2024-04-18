from playwright.sync_api import sync_playwright
from openai import OpenAI
from config import NEO4J_CLIENT
from helper import get_token_length, get_embedding, get_formatted_text


def Crawler(id: int, name: str) -> tuple[str, int]:

    info = NEO4J_CLIENT.get_user_info(name)
    if info is None:
        raise ValueError(f"User {name} not found")

    if info.openAIKey is None:
        raise ValueError(f"User {name} does not have an openAIKey")

    client = OpenAI(
        api_key=info.openAIKey,
    )

    repo = NEO4J_CLIENT.get_repo_info(id)

    if repo is None:
        raise ValueError(f"Repo with id {id} not found")

    if repo.html_url is None:
        raise ValueError(f"Repo with id {id} does not have an html_url")

    if repo.repo_vector is not None:
        return f"already generate embedding on repo: {id}", 201

    if repo.readme_summary is None or repo.readme_summary == "":
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

                repo.readme_summary = summary

                summary = summary.replace("\n", " ")
                vector = get_embedding(summary)
                repo.repo_vector = vector

                try:
                    NEO4J_CLIENT.save_repo_info(
                        repo_id=id,
                        readme_summary=repo.readme_summary,
                        repo_vector=repo.repo_vector,
                    )
                except Exception as e:
                    raise ValueError(
                        f"Error while inserting to neo4j: {str(id)}, {str(e)}"
                    )

                return f"success generate embedding on repo: {id}", 200

            except Exception as e:
                raise ValueError(f"Error while crawling: {str(e)}")

    else:
        summary = repo.readme_summary
        summary = summary.replace("\n", " ")

        vector = get_embedding(summary)

        repo.elk_vector = vector

        try:
            NEO4J_CLIENT.save_repo_info(
                repo_id=id,
                readme_summary=repo.readme_summary,
                repo_vector=repo.repo_vector,
            )

        except Exception as e:
            raise ValueError(f"Error while inserting to elastic: {str(id)}, {str(e)}")

        return f"success generate embedding on repo: {id}", 200


def Responser(name: str, text: str) -> list[dict]:

    info = NEO4J_CLIENT.get_user_info(name)
    if info is None:
        raise ValueError(f"User {name} not found")

    q = get_formatted_text(text)
    q = text.replace("\n", " ")

    vector = get_embedding(q)
    repos = NEO4J_CLIENT.get_suggestion_repos(name, info.limit, info.cosine, vector)

    return repos, 200
