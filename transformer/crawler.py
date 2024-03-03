from playwright.sync_api import sync_playwright
from model import db, RepoEmbeddingInfo
from openai import OpenAI
from config import NEO4J_CLIENT
from sqlalchemy import select
from helper import get_token_length
from model import db, RepoEmbeddingInfo


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

    if repo.summary_embedding is not None:
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
                vector = (
                    client.embeddings.create(
                        input=[article.text_content()], model="text-embedding-3-small"
                    )
                    .data[0]
                    .embedding
                )
                repo.summary_embedding = vector

                # vector = [TaggedDocument(words=get_tokens(summary), tags=['repo'])]

                # model = Doc2Vec(vector, vector_size=150, window=2, min_count=1, workers=1)
                # doc_vectors = model.dv['repo']

                # vector = array(doc_vectors, dtype=float32)
                # repo.summary_embedding = vector.tolist()

                db.session.commit()

                return f"success generate embedding on repo: {id}", 200

            except Exception as e:
                raise ValueError(f"Error while crawling: {str(e)}")

    else:
        summary = repo.readme_summary
        summary = summary.replace("\n", " ")

        # vector = [TaggedDocument(words=get_tokens(summary), tags=['repo'])]

        # model = Doc2Vec(vector, vector_size=150, window=2, min_count=1, workers=1)
        # doc_vectors = model.dv['repo']

        # vector = array(doc_vectors, dtype=float32)
        # repo.summary_embedding = vector.tolist()

        vector = (
            client.embeddings.create(input=[summary], model="text-embedding-3-small")
            .data[0]
            .embedding
        )
        repo.summary_embedding = vector

        db.session.commit()

        return f"success generate embedding on repo: {id}", 200


def Responser(name: str, text: str) -> list[dict]:

    info = NEO4J_CLIENT.get_user_info(name)

    if info is None:
        raise ValueError(f"User {name} not found")

    client = OpenAI(
        api_key=info["openAIKey"],
    )

    content = """
    You have a vector database containing summarized vectors of GitHub repositories that users have starred.
    Now, a user has presented a question.
    Please help me refine the following description for a more accurate expression, reply in English:
    {question}
    """

    text = text.replace("\n", " ")
    completion = client.chat.completions.create(
        messages=[
            {
                "role": "user",
                "content": content.format(question=text),
            }
        ],
        model="gpt-3.5-turbo",
    )

    text = completion.choices[0].message.content

    vector = (
        client.embeddings.create(input=[text], model="text-embedding-3-small")
        .data[0]
        .embedding
    )

    # vector = [TaggedDocument(words=get_tokens(text), tags=['repo'])]

    # model = Doc2Vec(vector, vector_size=150, window=2, min_count=1, workers=1)
    # doc_vectors = model.dv['repo']

    # vector = array(doc_vectors, dtype=float32)
    # vector = vector.tolist()

    # print(vector)

    items = db.session.scalars(
        select(RepoEmbeddingInfo)
        .filter(1 - RepoEmbeddingInfo.summary_embedding.cosine_distance(vector) > 0)
        .order_by(RepoEmbeddingInfo.summary_embedding.cosine_distance(vector))
        .limit(info["limit"])
    )

    return [item._to_json() for item in items], 200
