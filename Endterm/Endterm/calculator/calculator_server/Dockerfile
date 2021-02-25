FROM golang:1.15

WORKDIR /go/src/Endterm/calculator/calculator_server

COPY . /go/src/Endterm/calculator/calculator_server

RUN go mod download

EXPOSE 50051/tcp

CMD go run server.go