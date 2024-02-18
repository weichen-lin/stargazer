from flask import Flask, request, jsonify
from crawler import Crawler
from model import RepoEmbeddingInfoSchema, db
from pydantic import ValidationError
import os
from functools import wraps

app = Flask(__name__)


DATABASE_URL = os.environ.get("DATABASE_URL")
VALID_TOKEN = os.environ.get("AUTHENTICATION_TOKEN")

app.config["SQLALCHEMY_DATABASE_URI"] = DATABASE_URL

db.init_app(app)

def requires_auth(func):
    @wraps(func)
    def decorated(*args, **kwargs):
        auth_header = request.headers.get('Authorization')
        if not auth_header or auth_header != f"Bearer {VALID_TOKEN}":
            return jsonify({"error": "Unauthorized"}), 401
        return func(*args, **kwargs)
    return decorated

@app.route("/")
def healthy_check():
    return "ok", 200


@app.route("/vectorize", methods=["POST"])
@requires_auth
def vectorize():
    if request.is_json:
        data = request.get_json()

        try:
            model = RepoEmbeddingInfoSchema(**data)
            result, status = Crawler(model.repo_id)

            return jsonify({"message": result}), status

        except ValidationError as e:
            return jsonify({"error": str(e)}), 400

        except Exception as e:
            return jsonify({"error": str(e)}), 404

    else:
        return jsonify({"error": "Request must be JSON"}), 400


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8000, debug=True)
