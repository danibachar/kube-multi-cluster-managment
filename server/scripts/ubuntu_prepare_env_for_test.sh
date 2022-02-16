#!/bin/bash
# Hoping you have kind, kubectl and docker, and jsonnet installed locally

# Installations if needed

## Install jsonnet
apt-get update && install jsonnet

## Submariner CLI
curl -Ls https://get.submariner.io | bash
export PATH=$PATH:~/.local/bin
echo export PATH=\$PATH:~/.local/bin >> ~/.profile

# Submariner 1

## Submariner on Kind (sping 2 clusters)

git clone https://github.com/submariner-io/submariner
cd submariner
make clusters
cd ..

# Prometheus
## Prepare Promethues
git clone https://github.com/danibachar/submariner-cheatsheet
cd submariner-cheatsheet/prometheus/install
jb init
jb install github.com/prometheus-operator/kube-prometheus/jsonnet/kube-prometheus@release-0.8
jb update
sudo chmod +x ./build.sh
./build.sh

cd ../../..

## Deploy Prometheus on all clusters

kubectl --kubeconfig submariner/output/kubeconfigs/kind-config-cluster1 apply -f submariner-cheatsheet/prometheus/install/manifests/setup
kubectl --kubeconfig submariner/output/kubeconfigs/kind-config-cluster2 apply -f submariner-cheatsheet/prometheus/install/manifests/setup
sleep 5
kubectl --kubeconfig submariner/output/kubeconfigs/kind-config-cluster1 apply -f submariner-cheatsheet/prometheus/install/manifests/
kubectl --kubeconfig submariner/output/kubeconfigs/kind-config-cluster2 apply -f submariner-cheatsheet/prometheus/install/manifests/

# Submariner 2 - deploy Submariner (only after promethues)
## Note - broker is cluster1

## Deploy cluster1 as the Broker cluster
subctl deploy-broker --kubeconfig submariner/output/kubeconfigs/kind-config-cluster1
## joins cluster1 and cluster2 into the mesh
subctl join --kubeconfig submariner/output/kubeconfigs/kind-config-cluster1 broker-info.subm --clusterid cluster1 --natt=false
subctl join --kubeconfig submariner/output/kubeconfigs/kind-config-cluster2 broker-info.subm --clusterid cluster2 --natt=false

# Submariner 3 - exporting multi cluster services
subctl export service --kubeconfig submariner/output/kubeconfigs/kind-config-cluster1 --namespace monitoring prometheus-k8s
subctl export service --kubeconfig submariner/output/kubeconfigs/kind-config-cluster2 --namespace monitoring prometheus-k8s

# Local services and installations

kubectl --kubeconfig submariner/output/kubeconfigs/kind-config-cluster1 apply -f kube-multi-cluster-managment/server/go-serviceimports/kube/app.yaml
subctl export service --kubeconfig submariner/output/kubeconfigs/kind-config-cluster1 --namespace default serviceimports-svc

sleep 2m

kubectl --kubeconfig submariner/output/kubeconfigs/kind-config-cluster1 apply -f kube-multi-cluster-managment/server/py-koss/kube/job.yaml
