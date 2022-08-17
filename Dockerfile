FROM golang:1.19 as build

ENV CGO_ENABLED=0
WORKDIR /usr/src/app

COPY go.mod go.sum .
RUN go mod download && go mod verify

COPY ./ .
RUN go build -v -o /usr/local/bin/exercism ./exercism

CMD ["/usr/local/bin/exercism"]
