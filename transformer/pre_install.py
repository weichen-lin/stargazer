# from nltk import download
import tiktoken

# download("punkt", download_dir="/usr/share/nltk_data")
# download("stopwords", download_dir="/usr/share/nltk_data")

# download("punkt")
# download("stopwords")
encoding = tiktoken.get_encoding("cl100k_base")
