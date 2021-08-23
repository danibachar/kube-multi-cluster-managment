FROM golang:1.16

# RUN apk add --no-cache git
# RUN apk add --no-cache curl 

# These are temporary hack until we learn the actual API

# Instal clustercli
RUN curl -L https://github.com/kubernetes-sigs/cluster-api/releases/download/v0.3.1/clusterctl-linux-amd64 -o clusterctl
RUN chmod +x ./clusterctl
RUN mv ./clusterctl /usr/local/bin/clusterctl
# Install kubectl
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
RUN chmod +x ./kubectl
RUN mv ./kubectl /usr/local/bin

## We create an /app directory within our
## image that will hold our application source
## files
RUN mkdir /app
## We copy everything in the root directory
## into our /app directory
COPY /server/cluster/ /app
# Copy firebase json
COPY kmcm-firebase-admin.json /app
## We specify that we now wish to execute 
## any further commands inside our /app
## directory
WORKDIR /app
## Add this go mod download command to pull in any dependencies
RUN go mod download
## we run go build to compile the binary
## executable of our Go program
RUN go build -o /cluster

EXPOSE 8080

CMD ["/cluster"]