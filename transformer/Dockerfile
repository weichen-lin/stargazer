FROM --platform=linux/amd64 asia-east1-docker.pkg.dev/stargazer-414509/stargazer-images/stargazer-transformer:base

WORKDIR /app

COPY secrets.env /app/.env

ENV PORT $PATH
EXPOSE 8080

CMD ["python", "app.py"]