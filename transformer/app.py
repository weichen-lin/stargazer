from flask import Flask, request, jsonify
from crawler import Crawler
from model import RepoEmbeddingInfoSchema, db
from pydantic import ValidationError
import os

app = Flask(__name__)


db_url = os.environ.get("DATABASE_URL")

app.config["SQLALCHEMY_DATABASE_URI"] = db_url

db.init_app(app)


@app.route("/")
def healthy_check():
    return "ok", 200


@app.route("/vectorize", methods=["POST"])
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
