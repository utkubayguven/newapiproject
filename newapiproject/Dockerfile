# Build stage
FROM golang:1.22.2 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

# Run stage
FROM debian:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY config/config.json config/config.json
COPY .env .env
EXPOSE 8080
CMD ["./main"]


