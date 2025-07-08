FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd/server-core
RUN go build -o /server-core

FROM alpine:3.22
WORKDIR /
COPY --from=builder /server-core /server-core
COPY .env .env
EXPOSE 8080
CMD ["/server-core"] 