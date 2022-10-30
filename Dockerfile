FROM golang:1.18-alpine as builder

WORKDIR /go/src/vm

COPY go.mod go.sum ./
RUN go mod download
RUN go get gotest.tools/gotestsum

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-s' -o /go/bin/api -x /go/src/vm/cmd/api

EXPOSE 4000

CMD ["/go/bin/api"]