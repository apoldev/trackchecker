FROM golang:1.21-alpine3.18 as builder

WORKDIR /app

COPY ../go.mod go.sum ./
RUN go mod download
ADD .. /app/

RUN GOOS=linux go build ./cmd/trackchecker

ENTRYPOINT ./trackchecker

