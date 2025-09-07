# Stage 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o test-backend ./main.go

# Stage 2: Run
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/test-backend .

EXPOSE 8090

CMD ["./test-backend"]
