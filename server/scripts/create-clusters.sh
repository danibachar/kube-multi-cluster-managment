#!/usr/bin/env bash

# Copyright 2018 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script handles the creation of multiple clusters using kind and the
# ability to create and configure an insecure container registry.

set -o errexit
set -o nounset
set -o pipefail

# shellcheck source=util.sh
source "${BASH_SOURCE%/*}/util.sh"
NUM_CLUSTERS="${NUM_CLUSTERS:-2}"
KIND_IMAGE="${KIND_IMAGE:-}"
KIND_TAG="${KIND_TAG:-v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6}"
OS="$(uname)"

# Building Kind Clusters
function create-clusters() {
  local num_clusters=${1}

  local image_arg=""
  if [[ "${KIND_IMAGE}" ]]; then
    image_arg="--image=${KIND_IMAGE}"
  elif [[ "${KIND_TAG}" ]]; then
    image_arg="--image=kindest/node:${KIND_TAG}"
  fi
  for i in $(seq "${num_clusters}"); do
    kind create cluster --name "cluster${i}" "${image_arg}"
    fixup-cluster "${i}"
    echo

  done

  echo "Waiting for clusters to be ready"
  check-clusters-ready "${num_clusters}"
}

function fixup-cluster() {
  local i=${1} # cluster num

  if [ "$OS" != "Darwin" ];then
    # Set container IP address as kube API endpoint in order for clusters to reach kube API servers in other clusters.
    local docker_ip
    docker_ip=$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "cluster${i}-control-plane")
    kubectl config set-cluster "kind-cluster${i}" --server="https://${docker_ip}:6443"
  fi

  # Simplify context name
  kubectl config rename-context "kind-cluster${i}" "cluster${i}"
}

function check-clusters-ready() {
  for i in $(seq "${1}"); do
    util::wait-for-condition 'ok' "kubectl --context cluster${i} get --raw=/healthz &> /dev/null" 120
  done
}



# Promethues: Setting Up
function install-prometheus-on-clusters() {
    git clone https://github.com/danibachar/submariner-cheatsheet
    cd submariner-cheatsheet/prometheus/install
    jb init
    jb install github.com/prometheus-operator/kube-prometheus/jsonnet/kube-prometheus@release-0.8
    jb update
    sudo chmod +x ./build.sh
    ./build.sh
    cd ../../..

    local num_clusters=${1}

    for i in $(seq "${num_clusters}"); do
        kubectl config use-context "cluster${i}"
        kubectl apply -f submariner-cheatsheet/prometheus/install/manifests/setup
        sleep 2
        kubectl apply -f submariner-cheatsheet/prometheus/install/manifests/
        echo
    done
}

# Promethues: Export Multi-Cluster 

function export-prometheus-on-clusters() {
    local num_clusters=${1}

    for i in $(seq "${num_clusters}"); do
        kubectl config use-context "cluster${i}"
        subctl export service prometheus-k8s --namespace monitoring 
        echo
    done
}

# Submariner: Deploy and install
function install-subctl() {
    curl -Ls https://get.submariner.io | bash
    export PATH=$PATH:~/.local/bin
    echo export PATH=\$PATH:~/.local/bin >> ~/.profile
}

function set-cluster1-as-broker() {
    kubectl config use-context cluster1
    subctl deploy-broker
}

function deploy-submariner() {
    install-subctl
    set-cluster1-as-broker
    
    local num_clusters=${1}

    for i in $(seq "${num_clusters}"); do
        kubectl config use-context "cluster${i}"
        subctl join broker-info.subm --clusterid "cluster${i}" --natt=false
        echo
    done
}

# KOSS: Deploy and export services

function deploy-koss-services-on-clusters() {
    # Deploy optimizer only on borker
    kubectl config use-context cluster1
    kubectl apply -f py-koss/kube/job.yaml
    # Deploy koss services
    local num_clusters=${1}

    for i in $(seq "${num_clusters}"); do
        kubectl config use-context "cluster${i}"
        # Service-Imports
        kubectl apply -f go-serviceimports/kube/app.yaml
        # Service-Exporter
        kubectl apply -f go-serviceexporter/kube/app.yaml
        echo
    done
}

# function deploy-service-import() {
    
# }


echo "Creating ${NUM_CLUSTERS} clusters"
create-clusters "${NUM_CLUSTERS}"
install-prometheus-on-clusters "${NUM_CLUSTERS}"
deploy-submariner "${NUM_CLUSTERS}"
export-prometheus-on-clusters "${NUM_CLUSTERS}"
deploy-koss-services-on-clusters "${NUM_CLUSTERS}"

echo "Complete"