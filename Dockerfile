FROM golang:1.22.0 as builder

WORKDIR /app

COPY . .

RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/avito

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]