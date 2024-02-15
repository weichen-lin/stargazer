FROM --platform=linux/amd64 python:3.10.5-slim

WORKDIR /app

COPY . .

RUN pip install poetry \
    && poetry config virtualenvs.create false \
    && poetry install --no-dev \
    && poetry run playwright install chromium \
    && poetry run playwright install-deps

RUN poetry run python pre_install.py

CMD ["poetry", "run", "flask", "run"]
