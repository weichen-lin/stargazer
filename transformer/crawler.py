from cleaner import clean_text, get_tokens
from playwright.sync_api import sync_playwright
from gensim.models.doc2vec import Doc2Vec, TaggedDocument
from numpy import array, float32
from model import db, RepoEmbeddingInfo
from openai import OpenAI
import os

OPENAI_API_KEY = os.environ.get("OPENAI_API_KEY")

client = OpenAI(
    api_key=OPENAI_API_KEY,
)


def Crawler(id: int) -> tuple[str, int]:

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

    with sync_playwright() as p:
        try:
            browser = p.chromium.launch(headless=True)
            page = browser.new_page()
            page.goto(repo.html_url)

            article = page.wait_for_selector("article.markdown-body")
            chat_completion = client.chat.completions.create(
                messages=[
                    {
                        "role": "user",
                        "content": f"The following is a repository from github. After reading the material, can you give me a summary? Reply in English\n{article.text_content()}",
                    }
                ],
                model="gpt-3.5-turbo",
            )

            summary = chat_completion.choices[0].message.content

            repo.readme_summary = summary
            summary = clean_text(summary)

            vector = [TaggedDocument(words=get_tokens(summary), tags=["repo"])]

            model = Doc2Vec(vector, vector_size=150, window=2, min_count=1, workers=1)
            doc_vectors = model.dv["repo"]

            vector = array(doc_vectors, dtype=float32)
            repo.summary_embedding = vector.tolist()

            db.session.commit()

            return f"success generate embedding on repo: {id}", 200

        except Exception as e:
            raise ValueError(f"Error while crawling: {str(e)}")
