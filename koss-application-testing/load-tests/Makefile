# Kind
kind-clean: 
	kind delete clusters cluster-1

kind-dep:  kind-cluster promethues-prepare dep-test-app expose-svc

kind-cluster:
	kind create cluster --config=../kind/cluster.yaml
	kubectl cluster-info --context kind-cluster-1

# GCP
gcp-clean:
	gcloud container clusters delete --quiet "microsvc-us" --zone us-central1-a

gcp-dep: gcp-cluster gcp-cluster-auth promethues-prepare dep-test-app dep-ingress

gcp-cluster:
	gcloud container clusters create "microsvc-us" --zone us-central1-a --machine-type "n1-standard-1" --num-nodes 3 --update-addons=HttpLoadBalancing=ENABLED

gcp-cluster-auth:
	gcloud container clusters get-credentials microsvc-us --zone="us-central1-a"

# General
promethues-prepare: 
	helm install stable prometheus-community/kube-prometheus-stack
	kubectl patch --type=merge clusterrole stable-kube-prometheus-sta-operator --patch-file prometheus-operator-endpointslice-clusterrole.yaml
	helm upgrade --install -f patch-file.yaml stable prometheus-community/kube-prometheus-stack

dep-test-app:
	kubectl apply -f load_test.yaml

del-test-app:
	kubectl delete -f load_test.yaml

dep-ingress:
	kubectl apply -f ingress.yaml

del-ingress:
	kubectl delete -f ingress.yaml

expose-svc:
	kubectl port-forward svc/stable-grafana 8081:80 &
	kubectl port-forward svc/stable-kube-prometheus-sta-prometheus 8082:9090 &
	kubectl port-forward svc/simple-svc 8080:80 &

run-tester:
	./tester

## Installation
helm-update-repos:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo update

linux-install-helm:
	curl https://baltocdn.com/helm/signing.asc | sudo apt-key add -
	sudo apt-get install apt-transport-https --yes
	echo "deb https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
	sudo apt-get update
	sudo apt-get install helm

brew-install-helm:
	brew install helm