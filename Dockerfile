FROM golang:1.25.5 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o /out/todo ./cmd/server

FROM ubuntu:latest

WORKDIR /app

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db

COPY --from=builder /out/todo ./todo
COPY web ./web

EXPOSE 7540

CMD ["./todo"]
