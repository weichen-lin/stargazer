from graph import Neo4jOperations
import os

VALID_TOKEN = os.environ.get("AUTHENTICATION_TOKEN")
NEO4J_URI = os.environ.get("NEO4J_URL")
NEO4J_PASSWORD = os.environ.get("NEO4J_PASSWORD")
TRANSFORMER_PORT = os.environ.get("TRANSFORMER_PORT")
IS_PRODUCTION = os.environ.get("IS_PRODUCTION", True)

db = Neo4jOperations(NEO4J_URI, "neo4j", NEO4J_PASSWORD)
