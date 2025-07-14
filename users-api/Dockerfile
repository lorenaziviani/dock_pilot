FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY main.go .
RUN go build -o users-api main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/users-api .
EXPOSE 8080
ENTRYPOINT ["./users-api"]