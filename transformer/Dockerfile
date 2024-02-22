FROM --platform=linux/amd64 mcr.microsoft.com/playwright/python:v1.40.0-jammy

WORKDIR /app

COPY . .

RUN pip install poetry \
    && poetry config virtualenvs.create false \
    && poetry install --no-dev \
    && poetry run playwright install chromium \
    && poetry run python pre_install.py

COPY secrets.env /app/.env

ENV PORT $PATH
EXPOSE 8080

CMD ["python", "app.py"]