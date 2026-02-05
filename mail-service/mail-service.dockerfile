FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod  ./
RUN go mod download
COPY . .
RUN go build -o mail-service ./cmd/api


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/mail-service .
COPY --from=builder /app/templates ./templates
EXPOSE 6002
CMD ["./mail-service"]
