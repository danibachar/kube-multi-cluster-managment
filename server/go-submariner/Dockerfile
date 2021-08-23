FROM golang:1.16

# For Alpine
# RUN apk add --update \
#  python \
#  curl \
#  which \
#  bash

RUN apt-get update && apt-get install xz-utils

# Downloading gcloud package
RUN curl https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.tar.gz > /tmp/google-cloud-sdk.tar.gz
# Installing the package
RUN mkdir -p /usr/local/gcloud \
    && tar -C /usr/local/gcloud -xvf /tmp/google-cloud-sdk.tar.gz \
    && /usr/local/gcloud/google-cloud-sdk/install.sh --quiet
# Adding the package path to local
ENV PATH $PATH:/usr/local/gcloud/google-cloud-sdk/bin
RUN gcloud components install kubectl -q
## We create an /app directory within our
## image that will hold our application source
## files
RUN mkdir /app
## We copy everything in the root directory
## into our /app directory
COPY /server/go-submariner/ /app
COPY /server/golang-pkg/ /app/pkg/
# Copy firebase json
COPY kmcm-firebase-admin.json /app
COPY kmcm-owner.json /app
## We specify that we now wish to execute
## any further commands inside our /app
## directory
WORKDIR /app

# Auth GCLOUD
RUN gcloud auth activate-service-account --key-file=kmcm-owner.json

# Submariner CLI
RUN curl -Ls https://get.submariner.io | bash
RUN export PATH=$PATH:~/.local/bin
RUN echo export PATH=\$PATH:~/.local/bin >> ~/.profile
# RUN go install github.com/submariner-io/submariner-operator/pkg/subctl
# RUN mv ./kubectl /usr/local/bin
## Add this go mod download command to pull in any dependencies
RUN go mod download
## we run go build to compile the binary
## executable of our Go program
RUN go build -o /go-submariner

ENV SERVER_PORT=8080
ENV GOOGLE_APPLICATION_CREDENTIALS="./kmcm-owner.json"
EXPOSE 8080

CMD ["/go-submariner"]
