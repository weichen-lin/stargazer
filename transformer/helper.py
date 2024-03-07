from pre_install import encoding

# from sentence_transformers import SentenceTransformer


# string = 'I want a pizza'


# sentences = [string]

# model = SentenceTransformer('sentence-transformers/all-MiniLM-L6-v2')
# embeddings = model.encode(sentences)
# print(list(embeddings))


def get_token_length(test: str) -> int:
    return len(encoding.encode(test))


