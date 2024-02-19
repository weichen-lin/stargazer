FROM asia-east1-docker.pkg.dev/stargazer-414509/stargazer-images/stargazer-transformer:base

ENV AUTHENTICATION_TOKEN=$AUTHENTICATION_TOKEN
ENV DATABASE_URL=$DATABASE_URL
ENV OPENAI_API_KEY=$OPENAI_API_KEY

WORKDIR /app

EXPOSE 8000

CMD ["poetry", "run", "python", "app.py"]