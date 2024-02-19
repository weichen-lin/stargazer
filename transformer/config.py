from dotenv import load_dotenv, find_dotenv
import os

load_dotenv(find_dotenv())

DATABASE_URL = os.environ.get("DATABASE_URL")
VALID_TOKEN = os.environ.get("AUTHENTICATION_TOKEN")
OPENAI_API_KEY = os.environ.get("OPENAI_API_KEY")