# Multi-stage dockerfile
FROM golang:1.23 AS builder
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myApp .

FROM alpine:3.20
WORKDIR /server

# Install tzdata for time zone support
RUN apk add --no-cache tzdata

COPY --from=builder /app/myApp .
EXPOSE 8080
CMD ["./myApp"]