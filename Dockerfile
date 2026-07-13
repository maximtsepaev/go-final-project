FROM golang:1.25.5-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o /out/todo ./cmd/server

FROM alpine:3.20

WORKDIR /app

ENV TODO_PORT=7540

COPY --from=builder /out/todo ./todo
COPY web ./web

CMD ["./todo"]
