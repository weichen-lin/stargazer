import re
from nltk.stem.porter import PorterStemmer
from nltk.corpus import stopwords
from nltk.tokenize import word_tokenize
from model import RepoInfo

ps = PorterStemmer()


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


def flatten_and_deduplicate(nested_list):
    flat_list = []
    seen = set()

    def flatten(lst: list[RepoInfo]):
        for item in lst:
            if isinstance(item, list):
                flatten(item)
            else:
                if item.repo_id not in seen:
                    flat_list.append(item)
                    seen.add(item.repo_id)

    flatten(nested_list)
    return flat_list
