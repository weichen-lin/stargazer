import re
from nltk import download
from nltk.stem.porter import PorterStemmer
from nltk.corpus import stopwords
from nltk.tokenize import word_tokenize

ps = PorterStemmer()

download('punkt')
download('stopwords')

def clean_text(text):
    # 移除所有符號類型
    cleaned_text = re.sub(r"[^\w\s]", "", text)
    # 移除所有佔位符號，以空格代替
    cleaned_text = re.sub(r"\s+", " ", cleaned_text)
    # 去除頭尾空白
    cleaned_text = cleaned_text.strip()
    return cleaned_text

def get_tokens(text):
    # 將文本轉換成小寫
    text = text.lower()
    # 分詞
    tokens = word_tokenize(text)
    # 移除停用詞
    tokens = [word for word in tokens if word not in stopwords.words("english")]
    # 詞幹提取
    tokens = [ps.stem(word) for word in tokens]
    return tokens



# from gensim.models.doc2vec import Doc2Vec, TaggedDocument
# from nltk.tokenize import word_tokenize

# # 準備文檔數據
# documents = [TaggedDocument(words=word_tokenize(doc), tags=[str(i)]) for i, doc in enumerate(corpus)]

# # 訓練Doc2Vec模型
# model = Doc2Vec(documents, vector_size=50, window=2, min_count=1, workers=4)

# # 獲取文檔向量
# doc_vectors = [model.docvecs[str(i)] for i in range(len(corpus))]
# print(doc_vectors)
