#!/usr/bin/env bash

echo "Spinning up cluster"
kind create cluster --config=../kind/cluster.yaml
echo "Setting up cluster context"
kubectl cluster-info --context kind-cluster-1
echo "Deploy monitoring"
helm install stable prometheus-community/kube-prometheus-stack
echo "Deploy app"
kubectl apply -f load_test.yaml
echo "patching prometheus using helm to allow service discovery of our service (simple-svc"
kubectl patch --type=merge clusterrole stable-kube-prometheus-sta-operator --patch-file prometheus-operator-endpointslice-clusterrole.yaml
helm upgrade --install -f patch-file.yaml stable prometheus-community/kube-prometheus-stack

sleep 45

echo "Exposeing Graphana"
kubectl port-forward svc/stable-grafana 8081:80 &
echo "Exposeing Prometheus"
kubectl port-forward svc/stable-kube-prometheus-sta-prometheus 8082:9090 &
echo "Exposeing service locally"
kubectl port-forward svc/simple-svc 8080:80 &

echo "Patching Prometheus"

echo "Lunching tester"
# autocannon -R 1 -d 3600 -r 100 -c 20 -w 20 -m POST -H "Content-Type: application/json" -b '{"memory_params": {"duration_seconds": 0.2, "kb_count": 50}, "cpu_params": {"duration_seconds": 0.2, "load": 0.2}}' http://localhost:8080/load
# autocannon -R 8 -d 3600 -r 100 -c 20 -w 100 -m POST -H "Content-Type: application/json" -b '{"memory_params": {"duration_seconds": 0.2, "kb_count": 50}, "cpu_params": {"duration_seconds": 0.2, "load": 0.2}}' http://localhost:8080/load
