FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod  ./
RUN go mod download
COPY . .
RUN go build -o frontend ./cmd/web


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/frontend .
EXPOSE 8080
CMD ["./frontend"]
