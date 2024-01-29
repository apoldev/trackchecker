FROM golang:1.21-alpine3.18 as builder

RUN go install github.com/githubnemo/CompileDaemon@latest

WORKDIR /app

COPY ./go.mod ./
RUN go mod download
ADD . /app/

RUN GOOS=linux go build ./cmd/trackchecker

ENTRYPOINT CompileDaemon -directory=/app -build="go build ./cmd/trackchecker" -command=./trackchecker