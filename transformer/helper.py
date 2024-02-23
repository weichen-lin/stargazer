from openai import OpenAI
from config import OPENAI_API_KEY
from pre_install import encoding

client = OpenAI(
    api_key=OPENAI_API_KEY,
)

def get_embedding(text, model="text-embedding-3-small"):
   text = text.replace("\n", " ")
   return client.embeddings.create(input = [text], model=model).data[0].embedding

def get_token_length(test: str) -> int:
    return len(encoding.encode(test))