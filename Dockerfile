# multi-stage build
# Build Stage 
FROM arm64v8/golang:1.19-alpine3.17 AS builder

WORKDIR /app

COPY . .

RUN go build -o main main.go

# Run Stage
FROM arm64v8/alpine:3.17

WORKDIR /app

COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080

CMD ["/app/main"]