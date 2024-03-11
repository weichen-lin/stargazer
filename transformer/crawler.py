from playwright.sync_api import sync_playwright
from model import db, RepoEmbeddingInfo
from openai import OpenAI
from config import NEO4J_CLIENT
from helper import get_token_length, get_embedding, get_formatted_text
from elastic import insert_data, knn_search

def Crawler(id: int, name: str) -> tuple[str, int]:

    info = NEO4J_CLIENT.get_user_info(name)

    if info is None:
        raise ValueError(f"User {name} not found")

    client = OpenAI(
        api_key=info["openAIKey"],
    )

    repo = (
        db.session.query(RepoEmbeddingInfo)
        .filter(RepoEmbeddingInfo.repo_id == id)
        .first()
    )

    if repo is None:
        raise ValueError(f"Repo with id {id} not found")

    if repo.html_url is None:
        raise ValueError(f"Repo with id {id} does not have an html_url")

    if repo.elk_vector is not None:
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
                repo.elk_vector = vector

                try:
                    insert_data(id, repo._to_elastic())
                except Exception as e:
                    raise ValueError(f"Error while inserting to elastic: {str(id)}, {str(e)}")

                db.session.commit()

                return f"success generate embedding on repo: {id}", 200

            except Exception as e:
                raise ValueError(f"Error while crawling: {str(e)}")

    else:
        summary = repo.readme_summary
        summary = summary.replace("\n", " ")

        vector = get_embedding(summary)

        repo.elk_vector = vector

        db.session.commit()

        try:
            insert_data(id, repo._to_elastic())
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

    response = knn_search(vector)
    
    items = []

    if response['hits']['total']['value'] > 0:
        for hit in response['hits']['hits']:
            source = hit['_source']
            items.append({
                "full_name": source['full_name'],
                "avatar_url": source['avatar_url'],
                "html_url": source['html_url'],
                "description": source['description'],
                "stargazers_count": source['stargazers_count']
            })
        
        return items, 200
            
    else:
        return [], 200