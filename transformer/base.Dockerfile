FROM --platform=linux/amd64 mcr.microsoft.com/playwright/python:v1.40.0-jammy

WORKDIR /app

COPY . .

RUN pip install poetry \
    && poetry config virtualenvs.create false \
    && poetry install --no-dev \
    && poetry run playwright install chromium 

RUN poetry run python pre_install.py