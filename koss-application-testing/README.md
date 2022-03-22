# Kuberneets service benchmarking for load

## Prerequisits

On your management host (must have access to the Kubernetes cluster), you should install the following tools

- Homebrew - all of our installations are using homebrew, which is a package for macOS/linux (`https://stackoverflow.com/questions/33353618/can-i-use-homebrew-on-ubuntu/56982151`)
- Helm - used as a package manager for Kubernetes (using homebrew `brew install helm` / `arch -arm64 brew install helm`) 
  - Add repos:
    - `helm repo add bitnami https://charts.bitnami.com/bitnami`
    - `helm repo add prometheus-community https://prometheus-community.github.io/helm-charts`
  - Update repo
    - `helm repo update`

  - helm repo update 
- A clsuter (see link)
  - Cluster on the cloud
    - A Cloud account (GCP/WAS/Azure), for local testing using a kind cluster see - TODO - link
    - kubectl (`brew install kubernetes-cli`)
  - Local kind cluter
    - Docker (`brew install docker`)
    - kubectl (`brew install kubernetes-cli`)
    - kind (`brew install kind`)

## Cluster setup

### Kind Cluster

- From withing the kind folder run `kind create cluster --config=cluster.yaml`
- Set `kubectl` context for the cluster - `kubectl cluster-info --context kind-cluster-1`


### GCP - GKE Cluster

We are teesting on a small, hamble cluster

- Machine type, `t2d-standard-1` - 1 vCPUs, 4 GB RAM, SSD
- Cluster size, 3 machines
  
```
gcloud container clusters create "cluster-a"
--region "us-west1"
--machine-type "t2d-standard-1"
--image-type "UBUNTU"
--disk-type "pd-ssd"
--disk-size "15"
--num-nodes "3"
--no-enable-shielded-nodes
--no-shielded-integrity-monitoring
--no-shielded-secure-boot
--cluster-version "1.18.20-gke.4100"
--enable-ip-alias
--enable-network-policy
--enable-intra-node-visibility
--project=ivory-vim-337307
```

## After cluster setup installations

- Install Promethues
  - `helm install stable prometheus-community/kube-prometheus-stack`
- Now Promethues is up and running you can run the following commands to access it
  - `export POD_NAME=$(kubectl get pods --namespace default -l "app=prometheus,component=server" -o jsonpath="{.items[0].metadata.name}")`
  - `kubectl port-forward svc/stable-kube-prometheus-sta-prometheus 8082:9090`
  - kubect port-forward $POD_NAME 9090`
- You can access the graphana server as follow
  - `kubectl port-forward svc/stable-grafana 8080:3000`
  - The default user and password are [SO link](https://stackoverflow.com/questions/54039604/what-is-the-default-username-and-password-for-grafana-login-page)
    - user: `admin`
    - pwd: `prom-operator`


# Application / Test app deployment

- Deploy the load testing applicaition to see how it works under different loads
  - Deploy - `kubectl apply -f load_test.yaml`
  - Expose service - `kubectl port-forward svc/simple-svc 8080:80`
- Start load testing!
  - Install autocannon `npm i autocannon -g`


