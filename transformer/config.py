from dotenv import load_dotenv
from graph import UserInfo
import os

load_dotenv(dotenv_path="secrets.env")

DATABASE_URL = os.environ.get("DATABASE_URL")
VALID_TOKEN = os.environ.get("AUTHENTICATION_TOKEN")
NEO4J_URI = os.environ.get("NEO4J_URL")
NEO4J_PASSWORD = os.environ.get("NEO4J_PASSWORD")

NEO4J_CLIENT = UserInfo(NEO4J_URI, "neo4j", NEO4J_PASSWORD)