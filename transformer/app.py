from flask import Flask
from crawler import Crawler

app = Flask(__name__)

@app.route("/")
def hello():
    try:
        Crawler()
        return 'ok', 200
    except Exception as e:
        return str(e), 500



if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)