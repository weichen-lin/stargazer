from dotenv import load_dotenv
import os

load_dotenv(dotenv_path="secrets.env")

DATABASE_URL = os.environ.get("DATABASE_URL")
VALID_TOKEN = os.environ.get("AUTHENTICATION_TOKEN")
OPENAI_API_KEY = os.environ.get("OPENAI_API_KEY")