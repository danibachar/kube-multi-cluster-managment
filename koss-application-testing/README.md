# Kuberneets service benchmarking for load

## Prerequisits

On your management host (must have access to the Kubernetes cluster), you should install the following tools

- Homebrew - all of our installations are using homebrew, which is a package for macOS/linux (`https://stackoverflow.com/questions/33353618/can-i-use-homebrew-on-ubuntu/56982151`) 
- Helm - used as a package manager for Kubernetes (using homebrew `brew install helm` / `arch -arm64 brew install helm`) or for other linux distributions (`https://helm.sh/docs/intro/install/`)
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

We are testing on a small, hamble cluster 

- Machine type, `t2d-standard-1` - 1 vCPUs, 4 GB RAM, SSD
- Cluster size, 3 machines

It is recommended to run the commands from the GCP console, look at the README.md file for full automation

```
gcloud container clusters create "microsvc-us" \
--zone us-central1-a \
--machine-type "n1-standard-1" \
--num-nodes 3 \
--update-addons=HttpLoadBalancing=ENABLED
```

Get access to cluster

```
gcloud container clusters get-credentials microsvc-us --zone="us-central1-a"
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


## Application / Test app deployment

- Deploy the load testing applicaition to see how it works under different loads
  - Deploy - `kubectl apply -f load_test.yaml`
  - Expose service - `kubectl port-forward svc/simple-svc 8080:80`
- Start load testing!
  - Install autocannon `npm i autocannon -g`
- lunch attack `autocannon -R 1 -d 60 -r 5 -c 5 -w 5 -m POST -H "Content-Type: application/json" -b '{"memory_params": {"duration_seconds": 0.2, "kb_count": 50}, "cpu_params": {"duration_seconds": 0.2, "load": 0.2}}' http://localhost:8080/load`

##