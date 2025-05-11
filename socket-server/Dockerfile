# api1/Dockerfile (repeat for other APIs)
FROM golang:1.24-alpine

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source files
COPY . .

RUN go build -o socket-server main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/socket-server .

EXPOSE 10081

CMD ["./socket-server"]
