FROM golang:1.15

WORKDIR /go/src/Endterm/calculator/calculator_client

COPY . /go/src/Endterm/calculator/calculator_client

RUN apt-get update --yes
RUN apt-get install --yes netcat
RUN go mod download

EXPOSE 50051/tcp