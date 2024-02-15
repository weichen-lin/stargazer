from cleaner import clean_text, get_tokens
from playwright.sync_api import sync_playwright
from transformers import pipeline
from gensim.models.doc2vec import Doc2Vec, TaggedDocument

pipe = pipeline("summarization", model="Falconsai/text_summarization")

def Crawler():
    # 初始化Playwright
    with sync_playwright() as p:
        # 啟動瀏覽器
        browser = p.chromium.launch(headless=True)
        # 創建新頁面
        page = browser.new_page()
        # 前往網址
        page.goto("https://github.com/dubinc/dub")
        # 等待article標籤並且class名稱包含"markdown-body"的區塊出現
        article = page.wait_for_selector("article.markdown-body")
        # 獲取區塊的文本內容
        # article_text = article.text_content()
        summary = pipe(article.text_content(), max_length=300, min_length=30, do_sample=False)
        summary = summary[0]['summary_text']
        # print(summary)
        summary = clean_text(summary)
        
        vector = [TaggedDocument(words=get_tokens(summary), tags=['repo'])]
        
        model = Doc2Vec(vector, vector_size=50, window=2, min_count=1, workers=1)
        doc_vectors = [model.dv['repo']]
        print(doc_vectors)

        browser.close()
