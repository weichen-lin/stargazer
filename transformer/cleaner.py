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