# Multi-stage dockerfile
FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myApp .


FROM alpine:latest
WORKDIR /server
COPY --from=builder /app/myApp .
COPY --from=builder /app/templates ./templates
EXPOSE 80
CMD ["./myApp"]