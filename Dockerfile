FROM golang:1.24.1-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o appContract ./cmd/main.go

FROM alpine:latest
RUN apk add --no-cache postgresql-client
COPY --from=builder /app/appContract /appContract
EXPOSE 8080
CMD ["/appContract"]