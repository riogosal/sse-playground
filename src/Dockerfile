FROM golang:alpine

WORKDIR /go/src/app

COPY . .
RUN go mod tidy

ENTRYPOINT [ "go", "run", "main.go" ]