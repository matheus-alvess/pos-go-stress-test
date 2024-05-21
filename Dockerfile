FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go build -o stress-test .

ENTRYPOINT ["./stress-test"]
