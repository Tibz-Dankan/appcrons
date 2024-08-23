FROM golang:1.23.0-alpine

WORKDIR /app

COPY . .

RUN go build -o ./bin/appcrons ./cmd

EXPOSE 8080

CMD ["./bin/appcrons"]
