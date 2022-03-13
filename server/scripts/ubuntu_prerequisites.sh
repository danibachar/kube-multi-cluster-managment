#!/bin/bash
# Hoping you have kind, kubectl and docker, and jsonnet installed locally

# Installations if needed

## Make sure minimal go version installed
echo "installing go 1.13"
sudo apt-get purge golang* -y
sudo rm -rvf /usr/local/go

wget https://go.dev/dl/go1.13.15.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.13.15.linux-amd64.tar.gz

mkdir ~/.go

echo export GOROOT=/usr/local/go >> ~/.profile
echo export GOPATH=$HOME/.go >> ~/.profile

echo export PATH=$PATH:$(go env GOPATH)/bin >> ~/.profile
echo export PATH=$GOPATH/bin:$GOROOT/bin:$PATH >> ~/.profile

echo export GO111MODULE=on >> ~/.profile

source ~/.profile

echo export GOROOT=/usr/local/go >> ~/.bashrc
echo export GOPATH=$HOME/.go >> ~/.bashrc

echo export PATH=$PATH:$(go env GOPATH)/bin >> ~/.bashrc
echo export PATH=$GOPATH/bin:$GOROOT/bin:$PATH >> ~/.bashrc

echo export GO111MODULE=on >> ~/.bashrc

source ~/.bashrc

sudo update-alternatives --install "/usr/bin/go" "go" "/usr/local/go/bin/go"
sudo update-alternatives --set go /usr/local/go/bin/go

echo "Installing jsonnet + jsonnet-bundler"
go get -a github.com/jsonnet-bundler/jsonnet-bundler/cmd/jb
go get -a github.com/brancz/gojsontoyaml@latest
go get -a github.com/google/go-jsonnet/cmd/jsonnet
# OR
brew install go-jsonnet

echo "Installing Submariner CLI"
curl -Ls https://get.submariner.io | bash
export PATH=$PATH:~/.local/bin
echo export PATH=\$PATH:~/.local/bin >> ~/.profilek

echo "Allow sch_netem to simulate network delay between cluster nodes"
sudo modprobe sch_netem