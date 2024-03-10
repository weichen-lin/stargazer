from elasticsearch import Elasticsearch
from model import ElasticSearchDoc
from config import ELASTIC_PASSWORD, ELASTIC_CLOUD_ID


# Create the client instance
client = Elasticsearch(
    cloud_id=ELASTIC_CLOUD_ID,
    basic_auth=("elastic", ELASTIC_PASSWORD)
)

def insert_data(id: str, data :ElasticSearchDoc):
    client.index(index="repository", id=id, body=data)

def knn_search(vector: list):
    query_body = {
        "knn": {
            "field": "elk_vector",
            "query_vector": vector,
            "k": 10,
            "num_candidates": 100
        },
        "_source": {
            "includes": ["full_name", "avatar_url", "html_url", "description", "stargazers_count"]
        }
    }

    res = client.search(index="repository", body=query_body)
    
    return res