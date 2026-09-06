# Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git
WORKDIR /app

COPY go.mod ./
COPY . .
ENV GOFLAGS=-mod=mod

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate ./cmd/migrate

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates wget
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/migrate .

EXPOSE 8003
CMD ["./main"]
