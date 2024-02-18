FROM --platform=linux/amd64 golang:1.21.0-alpine AS builder

WORKDIR /app

ENV POSTGRESQL_URL=postgresql://johndoe:randompassword@stargazer-postgres:5432/mydb?sslmode=disable
ENV NEO4J_URL=neo4j://stargazer-neo4j:7687
ENV NEO4J_PASSWORD=randompassword

COPY . .

RUN go mod download