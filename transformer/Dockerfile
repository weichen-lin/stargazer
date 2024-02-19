FROM asia-east1-docker.pkg.dev/stargazer-414509/stargazer-images/stargazer-transformer:base

WORKDIR /app

COPY secrets.env /app/.env

EXPOSE 8000

CMD ["poetry", "run", "python", "app.py"]