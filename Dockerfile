# First stage: Build the binary
FROM golang:1.24.1-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . . 
RUN go build -ldflags="-s -w" -o /main cmd/main.go

# Second stage: Minimal final image
FROM alpine:3.21
WORKDIR /root/

COPY --from=builder /main .

CMD ["./main"]
