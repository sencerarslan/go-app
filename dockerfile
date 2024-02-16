FROM golang:latest

WORKDIR /go-app

COPY . .

RUN go mod download

RUN go build -o main .

CMD ["./main"]
