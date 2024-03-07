from elasticsearch import Elasticsearch


# client = Elasticsearch(
#     "https://localhost:9200",  # Elasticsearch endpoint
#     basic_auth=("elastic", "test_password"),  # Elasticsearch credentials
#     verify_certs=False,  # Disable SSL certificate verification
#     ca_certs=False,  # Disable CA certificates
# )

client = Elasticsearch(
    "https://localhost:9200",  # Elasticsearch endpoint
    api_key=('for_local_test', '6g4PdlmUSo-AIVSQSYxxwA'),  # API key ID and secret
)

client.indices.create(index="my_index")
