# from nltk import download
import tiktoken
from sentence_transformers import SentenceTransformer
from transformers import T5Tokenizer, T5ForConditionalGeneration

model = SentenceTransformer("sentence-transformers/all-MiniLM-L6-v2")

# download("punkt", download_dir="/usr/share/nltk_data")
# download("stopwords", download_dir="/usr/share/nltk_data")

encoding = tiktoken.get_encoding("cl100k_base")

tokenizer = T5Tokenizer.from_pretrained("google/flan-t5-base", legacy=True)
t5_model = T5ForConditionalGeneration.from_pretrained("google/flan-t5-base")
