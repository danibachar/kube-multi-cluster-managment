FROM golang:1.16

RUN mkdir /app

## We copy everything in the root directory
## into our /app directory
COPY /server/go-serviceimports/ /app
## We specify that we now wish to execute
## any further commands inside our /app
## directory
WORKDIR /app
## Add this go mod download command to pull in any dependencies
RUN go mod download
## we run go build to compile the binary
## executable of our Go program
RUN go build -o /go-serviceexporter

ENV SERVER_PORT=8080
EXPOSE 8080

CMD ["/go-serviceexporter"]
