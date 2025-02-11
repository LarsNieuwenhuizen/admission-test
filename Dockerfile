FROM golang:1.23-alpine AS builder
LABEL authors="Lars Nieuwenhuizen"

RUN mkdir /app

WORKDIR /app

COPY . /app

RUN mkdir bin; \
    go mod tidy

RUN go build -o bin/webhook main.go

FROM alpine:3.21 AS result
LABEL authors="Lars Nieuwenhuizen"

COPY --from=builder /app/bin/webhook /app/bin/webhook

ENTRYPOINT ["/app/bin/webhook"]
CMD ["run"]