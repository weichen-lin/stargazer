from dotenv import load_dotenv
from graph import Neo4jOperations
import os

load_dotenv(dotenv_path="secrets.env")

VALID_TOKEN = os.environ.get("AUTHENTICATION_TOKEN")
NEO4J_URI = os.environ.get("NEO4J_URL")
NEO4J_PASSWORD = os.environ.get("NEO4J_PASSWORD")

NEO4J_CLIENT = Neo4jOperations(NEO4J_URI, "neo4j", NEO4J_PASSWORD)
