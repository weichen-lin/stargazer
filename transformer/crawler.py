from cleaner import clean_text, get_tokens
from playwright.sync_api import sync_playwright
from gensim.models.doc2vec import Doc2Vec, TaggedDocument
from pre_install import pipe
from numpy import array, float32
from model import db, RepoEmbeddingInfo
from sqlalchemy import cast, BigInteger

def Crawler(id: int):
    # return id
    repo = db.session.query(RepoEmbeddingInfo).filter(RepoEmbeddingInfo.repo_id == id).first()

    if repo is None:
        raise ValueError(f"Repo with id {id} not found")

    if repo.html_url is None:
        raise ValueError(f"Repo with id {id} does not have an html_url")

    with sync_playwright() as p:
        # 啟動瀏覽器
        try:
            browser = p.chromium.launch(headless=True)
            # 創建新頁面
            page = browser.new_page()
            # 前往網址
            page.goto(repo.html_url)

            article = page.wait_for_selector("article.markdown-body")

            summary = pipe(article.text_content(), max_length=300, min_length=30, do_sample=False)
            summary = summary[0]['summary_text']
            
            repo.readme_summary = summary
            # print(summary)
            summary = clean_text(summary)
            
            vector = [TaggedDocument(words=get_tokens(summary), tags=['repo'])]
            
            model = Doc2Vec(vector, vector_size=150, window=2, min_count=1, workers=1)
            doc_vectors = model.dv['repo']


            vector = array(doc_vectors, dtype=float32)
            repo.summary_embedding = vector.tolist()
            db.session.commit()
        
        except Exception as e:
            raise ValueError(f"Error while crawling: {str(e)}")

        finally:
            browser.close()
