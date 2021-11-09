####### Cluster
build-cluster:
	docker build -f go-cluster/Dockerfile -t danibachar/cluster-msvc:latest  ..

#&& docker push danibachar/cluster-msvc:latest
#--no-cache

run-cluster: build-cluster
	docker run -p 8080:8080 -it danibachar/cluster-msvc

deploy-cluster: build-cluster
	kubectl apply -f go-cluster/deployment.yaml

delete-cluster:
	kubectl delete -f go-cluster/deployment.yaml

expose-cluster:
	kubectl port-forward service/cluster-msvc 8080:8080

####### Poviders
build-providers:
	docker build -f go-providers/Dockerfile -t danibachar/providers-msvc:latest  ..

#&& docker push danibachar/cluster-msvc:latest
#--no-cache

run-providers: build-providers
	docker run -p 8081:8081 -it danibachar/providers-msvc

deploy-providers: build-providers
	kubectl apply -f go-providers/deployment.yaml

delete-providers:
	kubectl delete -f go-providers/deployment.yaml

expose-providers:
	kubectl port-forward service/providers-msvc 8081:8081

####### Submariner
build-submariner:
	docker build -f go-submariner/Dockerfile -t danibachar/submariner-msvc:latest  ..

run-submariner: build-submariner
	docker run -p 8080:8080 -it danibachar/submariner-msvc

deploy-submariner: build-submariner
	kubectl apply -f go-submariner/deployment.yaml

delete-submariner:
	kubectl delete -f go-submariner/deployment.yaml

expose-submariner:
	kubectl port-forward service/submariner-msvc 8080:8080


####### Service Imports
build-serviceimports:
	docker build -f go-serviceimports/Dockerfile -t danibachar/serviceimports-msvc:latest  ..

run-serviceimports:	build-serviceimports
	docker run -p 8080:8080 -it danibachar/serviceimports-msvc:latest

release-serviceimports:	build-serviceimports
	docker push danibachar/serviceimports-msvc:latest

deploy-serviceimports:	release-serviceimports
	kubectl apply -f go-serviceimports/kube/app.yaml

delete-serviceimports:
	kubectl delete -f go-serviceimports/kube/app.yaml

expose-serviceimports:
	kubectl port-forward service/serviceimports-msvc 8080:8080

####### Service Exports
build-serviceexporter:
	docker build -f go-serviceexporter/Dockerfile -t danibachar/serviceexporter-msvc:latest  ..

run-serviceexporter:	build-serviceimports
	docker run -p 8080:8080 -it danibachar/serviceexporter-msvc:latest

release-serviceexporter:	build-serviceexporter
	docker push danibachar/serviceexporter-msvc:latest

deploy-serviceexporter:	release-serviceexporter
	kubectl apply -f go-serviceexporter/kube/app.yaml

delete-serviceexporter:
	kubectl delete -f go-serviceexporter/kube/app.yaml

expose-serviceexporter:
	kubectl port-forward service/serviceexporter-msvc 8080:8080

####### Optimization Job - KOSS
build-koss:
	docker build -f py-koss/Dockerfile -t danibachar/koss-job:latest  ..

release-koss:	build-koss
	docker push danibachar/koss-job:latest

deploy-koss:	release-koss
	kubectl apply -f py-koss/kube/app.yaml

delete-koss:
	kubectl delete -f py-koss/kube/app.yaml

####### General

release-all:	release-serviceimports release-serviceexporter release-koss

get-all:
	kubectl get pods
	kubectl get deployment
	kubectl get services

describe-all:
	kubectl describe pods
	kubectl describe deployment
	kubectl describe services
