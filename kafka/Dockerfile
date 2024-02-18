FROM --platform=linux/amd64 golang:1.21.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main.exe

FROM --platform=linux/amd64 alpine:latest

WORKDIR /app

EXPOSE 8080

ENV POSTGRESQL_URL=postgresql://johndoe:randompassword@stargazer-postgres:5432/mydb?sslmode=disable
ENV NEO4J_URL=neo4j://stargazer-neo4j:7687
ENV NEO4J_PASSWORD=randompassword

COPY --from=builder /app/main.exe .