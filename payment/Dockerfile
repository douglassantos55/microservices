FROM golang:1.19

WORKDIR /app

COPY . .

RUN go mod download && \
    go build -v -o /usr/local/bin/app ./cmd

EXPOSE 80

CMD ["app"]
