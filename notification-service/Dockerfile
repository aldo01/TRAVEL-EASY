FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git
WORKDIR /app

COPY go.mod ./
COPY . .
RUN go mod tidy && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o notification-service .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/notification-service .

EXPOSE 8090
CMD ["./notification-service"]
