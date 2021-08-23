FROM golang:1.16

# For Alpine
# RUN apk add --update \
#  python \
#  curl \
#  which \
#  bash

RUN apt-get update && apt-get install xz-utils uuid-runtime

# Install kubectl
# RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
# RUN chmod +x ./kubectl
# RUN mv ./kubectl /usr/local/bin

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
COPY /server/go-providers/ /app
# COPY /server/golang-pkg/models /app/pkg/models
# COPY /server/golang-pkg/utils /app/pkg/utils
# Copy firebase json
COPY kmcm-firebase-admin.json /app
COPY kmcm-owner.json /app
## We specify that we now wish to execute
## any further commands inside our /app
## directory
WORKDIR /app

# Auth GCLOUD
RUN gcloud auth activate-service-account --key-file=kmcm-owner.json

# Allow GCP Configuration scripts
RUN wget https://raw.githubusercontent.com/sridhargaddam/k8sscripts/main/rp_filter_settings/update-rp-filter.sh
RUN wget https://raw.githubusercontent.com/sridhargaddam/k8sscripts/main/rp_filter_settings/configure-rp-filter.sh
RUN chmod +x update-rp-filter.sh
RUN chmod +x configure-rp-filter.sh
# RUN mv ./kubectl /usr/local/bin
## Add this go mod download command to pull in any dependencies
RUN go mod download
## we run go build to compile the binary
## executable of our Go program
RUN go build -o /go-providers

ENV SERVER_PORT=8081
ENV GOOGLE_APPLICATION_CREDENTIALS="./kmcm-owner.json"
EXPOSE 8081

CMD ["/go-providers"]
