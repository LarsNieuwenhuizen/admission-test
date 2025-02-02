FROM golang:1.23-alpine
LABEL authors="Lars Nieuwenhuizen"

RUN mkdir /app

WORKDIR /app

COPY . /app

RUN mkdir bin; \
    go mod tidy

RUN go build -o bin/webhook main.go

ENTRYPOINT ["/app/bin/webhook"]
CMD ["run"]