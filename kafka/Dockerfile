FROM --platform=linux/amd64 golang:1.21.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GIN_MODE=release

RUN go build -o main.exe

FROM --platform=linux/amd64 alpine:latest

WORKDIR /app

COPY secrets.env /app/secrets.env
COPY templates /app/templates

EXPOSE 8080

COPY --from=builder /app/main.exe .

CMD ["./main.exe"]