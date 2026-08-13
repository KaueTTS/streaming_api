FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN GOOS=linux go build -o streaming-api .

FROM alpine:3.21
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/streaming-api ./streaming-api
COPY --from=builder /app/.env.example ./.env
EXPOSE 8080
CMD ["./streaming-api"]
