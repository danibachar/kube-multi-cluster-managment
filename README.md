# kube-multi-cluster-managment
Utilizing Kubernetes Cluster-API, Multi-Cluster API and Submariner to build management and observation tool for researching the Kuberntes multi-cluster environment

## Installation
- Docker
- kubectl
- kind
- homebrew

## Providers
- AWS - needs Administrator access permissions

## Step by step - Kind
Based on - `https://cluster-api.sigs.k8s.io/user/quick-start.html`

Create clsuter with supported version
kind create cluster --image=kindest/node:v1.22.0

AWS
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=
export AWS_SECRET_ACCESS_KEY=
export AWS_B64ENCODED_CREDENTIALS=$(clusterawsadm bootstrap credentials encode-as-profile)

GCP
export GCP_B64ENCODED_CREDENTIALS=$( cat ./kmcm-owner.json | base64 | tr -d '\n' )

export GCP_REGION=us-east1 GCP_PROJECT=kmcm-83960 GCP_NODE_MACHINE_TYPE=n1-standard-2 GCP_NETWORK_NAME=default GCP_CONTROL_PLANE_MACHINE_TYPE=n1-standard-2

And create a cluster

clusterctl generate cluster aws-us-east-1 \
--kubernetes-version v1.22.0 \
--control-plane-machine-count=3 \
--worker-machine-count=3 \
| kubectl apply -f -

clusterctl config cluster test1 --kubernetes-version v1.18.16 --control-plane-machine-count=3 --worker-machine-count=3 | kubectl apply -f -


## Simple tesing on a local managment cluster - creating a work load cluster


## Docker cheatsheet

- Clean up local docker env
`docker rm -vf $(docker ps -a -q) && docker rmi -f $(docker images -a -q)`


## AWS Mahcine types
- `https://aws.amazon.com/ec2/instance-types/`

## Kubectl Cheetsheet
- `https://kubernetes.io/docs/reference/kubectl/cheatsheet/`


kubectl --kubeconfig=./capi-quickstart.kubeconfig apply -f https://docs.projectcalico.org/v3.15/manifests/calico.yaml

kubectl --kubeconfig=./capi-quickstart.kubeconfig apply -f https://docs.projectcalico.org/v3.18/manifests/calico.yaml


## GCP Allow default network

gcloud compute routers create capi-quickstart-myrouter --project=kmcm-83960 --region=us-east1 --network=default

gcloud compute routers nats create capi-quickstart --project=kmcm-83960 --router-region=us-east1 --router=capi-quickstart-myrouter --nat-all-subnet-ip-ranges --auto-allocate-nat-external-ips

gcloud compute firewall-rules list --project kmcm-83960

gcloud compute networks list --project=kmcm-83960

gcloud compute networks describe default --project=kmcm-83960

kubectl --kubeconfig=./capi-quickstart.kubeconfig get nodes



// Cluster in asia
gcloud container clusters create test-europe-west1 --region=europe-west1 --machine-type=e2-micro --num-nodes=1

gcloud container clusters get-credentials test-asia-east1 --region=asia-east1 --project=kmcm-83960
~/.local/bin/subctl deploy-broker
~/.local/bin/subctl join broker-info.subm --clusterid cluster-a --servicecidr 10.7.240.0/20

// Cluster in US
gcloud container clusters create test-us-east1 --region=us-east1 --machine-type=e2-micro --num-nodes=1

gcloud container clusters get-credentials test-us-east1 --region=us-east1 --project=kmcm-83960
~/.local/bin/subctl join broker-info.subm --clusterid cluster-a --servicecidr 10.3.240.0/20



KUBECONFIG=test-asia-east1.yml gcloud container clusters get-credentials test-asia-east1 --region=asia-east1 --project=kmcm-83960
KUBECONFIG=test-us-east1.yml gcloud container clusters get-credentials test-us-east1 --region=us-east1 --project=kmcm-83960

KUBECONFIG=test-us-east1.yml:test-asia-east1.yml ~/.local/bin/subctl verify --kubecontexts test-us-east1,test-asia-east1 --only service-discovery,connectivity --verbose

- Get gatway nodes
kubectl get nodes --selector='submariner.io/gateway=true' --all-namespaces

1)
gcloud container clusters create "cluster-a" \
    --region "us-west1" \
    --machine-type "g1-small" \
    --image-type "UBUNTU" \
    --disk-type "pd-ssd" \
    --disk-size "15" \
    --num-nodes "1" \
    --no-enable-shielded-nodes \
    --no-shielded-integrity-monitoring \
    --no-shielded-secure-boot \
    --cluster-version "1.18.20-gke.4100" \
    --enable-ip-alias \
    --enable-network-policy \
    --enable-intra-node-visibility \
    --project=kmcm-83960

gcloud container clusters create "cluster-b" \
    --region "us-east1" \
    --machine-type "g1-small" \
    --image-type "UBUNTU" \
    --disk-type "pd-ssd" \
    --disk-size "15" \
    --num-nodes "1" \
    --no-enable-shielded-nodes \
    --no-shielded-integrity-monitoring \
    --no-shielded-secure-boot \
    --cluster-version "1.18.20-gke.4100" \
    --enable-ip-alias \
    --enable-network-policy \
    --enable-intra-node-visibility \
    --project=kmcm-83960

gcloud container clusters delete cluster-a --region="us-west1"
gcloud container clusters delete cluster-b --region="us-east1"

2)
gcloud container clusters get-credentials cluster-a --region="us-west1"
./configure-rp-filter.sh
gcloud container clusters get-credentials cluster-b --region="us-east1"
./configure-rp-filter.sh

gcloud container clusters get-credentials cluster-a --zone="europe-west3-a"
./configure-rp-filter.sh
gcloud container clusters get-credentials cluster-b --zone="europe-west3-a"
./configure-rp-filter.sh

3)
gcloud compute firewall-rules create "allow-tcp-in" --allow=tcp \
  --direction=IN --source-ranges=10.12.0.0/20,10.8.0.0/14,10.4.0.0/20,10.0.0.0/14

gcloud compute firewall-rules create "allow-tcp-out" --allow=tcp --direction=OUT \
  --destination-ranges=10.12.0.0/20,10.8.0.0/14,10.4.0.0/20,10.0.0.0/14

gcloud compute firewall-rules create "udp-in-500" --allow=udp:500 --direction=IN
gcloud compute firewall-rules create "udp-in-4500" --allow=udp:4500 --direction=IN
gcloud compute firewall-rules create "udp-in-4800" --allow=udp:4800 --direction=IN

gcloud compute firewall-rules create "udp-out-500" --allow=udp:500 --direction=OUT
gcloud compute firewall-rules create "udp-out-4500" --allow=udp:4500 --direction=OUT
gcloud compute firewall-rules create "udp-out-4800" --allow=udp:4800 --direction=OUT

4)
gcloud container clusters get-credentials cluster-a --zone="europe-west3-a"
subctl deploy-broker

gcloud container clusters get-credentials cluster-a --region="us-west1"
subctl deploy-broker

5)

gcloud container clusters get-credentials cluster-a --zone=europe-west3-a --project=kmcm-83960
subctl join broker-info.subm --clusterid cluster-a --clustercidr 10.88.0.0/14 --servicecidr 10.92.0.0/20 --health-check=false

gcloud container clusters get-credentials cluster-b --zone=europe-west3-a --project=kmcm-83960
subctl join broker-info.subm --clusterid cluster-b --clustercidr 10.120.0.0/14 --servicecidr 10.124.0.0/20 --health-check=false

6)
gcloud container clusters get-credentials cluster-a --zone=europe-west3-a --project=kmcm-83960
subctl show all

7)
gcloud container clusters get-credentials cluster-a --zone=europe-west3-a --project=kmcm-83960
CLUSTER_IP=$(kubectl get svc submariner-lighthouse-coredns -n submariner-operator -o=custom-columns=ClusterIP:.spec.clusterIP | tail -n +2)


gcloud container clusters get-credentials cluster-b --zone=europe-west3-a --project=kmcm-83960
CLUSTER_IP=$(kubectl get svc submariner-lighthouse-coredns -n submariner-operator -o=custom-columns=ClusterIP:.spec.clusterIP | tail -n +2)


kubectl config delete-cluster gke_kmcm-83960_europe-west3-a_cluster-a
kubectl config delete-context gke_kmcm-83960_europe-west3-a_cluster-a
kubectl config delete-user gke_kmcm-83960_europe-west3-a_cluster-a

kubectl config delete-cluster gke_kmcm-83960_europe-west3-a_cluster-b
kubectl config delete-context gke_kmcm-83960_europe-west3-a_cluster-b
kubectl config delete-user gke_kmcm-83960_europe-west3-a_cluster-b

# E2E testing

KUBECONFIG=cluster-a.yml gcloud container clusters get-credentials cluster-a --zone="europe-west3-a"
KUBECONFIG=cluster-b.yml gcloud container clusters get-credentials cluster-b --zone="europe-west3-a"

KUBECONFIG=cluster-a.yml:cluster-b.yml subctl verify --kubecontexts gke_kmcm-83960_europe-west3-a_cluster-a,gke_kmcm-83960_europe-west3-a_cluster-b --only service-discovery,connectivity --verbose 

# Manual Testing

 #### Deploy on cluster-a
gcloud container clusters get-credentials cluster-a --zone=europe-west3-a --project=kmcm-83960
kubectl create deployment nginx --image=nginx
kubectl expose deployment nginx --port=80
subctl export service --namespace default nginx

#### Run test from cluster-b
gcloud container clusters get-credentials cluster-b --zone=europe-west3-a --project=kmcm-83960
kubectl -n default run tmp-shell --rm -i --tty --image quay.io/submariner/nettest -- /bin/bash
curl nginx.default.svc.clusterset.local

Cluster-a
gke-cluster-a-default-pool-29836d96-7c9x/10.156.0.30
10.108.0.14

Cluster-b


### Delete netshoot pods
kubectl get pods | grep netshoot-hostmount | awk '/netshoot-hostmoun/{print $1}' | xargs kubectl delete pod

kubectl -n default run tmp-shell --privileged --rm -i --tty --image nicolaka/netshoot -- /bin/bash

kubectl run --privileged netshoot-hostmount-$(uuidgen) -i --overrides='{
	"spec": {
		"hostNetwork": true,
		"nodeName": "gke-cluster-a-default-pool-573bbd6a-7skq",
		"containers": [{
			"stdin": true,
			"stdinOnce": true,
			"terminationMessagePath": "/dev/termination-log",
			"terminationMessagePolicy": "File",
			"tty": true,
			"securityContext": {
				"allowPrivilegeEscalation": true,
				"privileged": true,
				"runAsUser": 0,
				"capabilities": {
					"add": ["ALL"]
				}
			},
			"name": "netshoot-hostmount",
			"image": "nicolaka/netshoot",
			"volumeMounts": [{
				"mountPath": "/host",
				"name": "host-slash",
				"readOnly": true
			}]
		}],
	        "restartPolicy": "Never",
		"volumes": [{
			"hostPath": {
				"path": "/",
				"type": ""
			},
			"name": "host-slash"
		}]
	}
}' --image nicolaka/netshoot -- /bin/bash

sysctl -a 2>/dev/null | grep "\.rp_filter" | awk '/net.ipv4/{print $1}' | tr . / | awk '/net/{newvar="/proc/sys/"$1; print newvar}' | awk '{print "2" > $1}'


kubectl get pods | grep netshoot-hostmount | awk '/netshoot-hostmoun/{print $1}' | xargs kubectl delete pod

sysctl -a 2>/dev/null | grep '\.rp_filter' | awk '/net.ipv4/{print $1}' | tr . / | sudo awk '/net/{newvar="/proc/sys/"$1; print newvar}' | awk '{print "0" > $1}'

sudo 


echo 2 > /proc/sys/net/ipv4/conf/vx-submariner/rp_filter

echo 2 > /proc/sys/net/ipv4/conf/ens4/rp_filter && echo 2 > /proc/sys/net/ipv4/conf/docker0/rp_filter && echo 2 > /proc/sys/net/ipv4/conf/cbr0/rp_filter


## Calico config


KUBECONFIG=cluster-a.yml gcloud container clusters get-credentials cluster-a --region="us-west1"
KUBECONFIG=cluster-b.yml gcloud container clusters get-credentials cluster-b --region="us-east1"

curl -o kubectl-calico -O -L  "https://github.com/projectcalico/calicoctl/releases/download/v3.16.10/calicoctl" 
chmod +x kubectl-calico

gcloud container clusters get-credentials cluster-a --zone=europe-west3-a --project=kmcm-83960

- Cluster a

cat > svcclusterb.yaml <<EOF
  apiVersion: projectcalico.org/v3
  kind: IPPool
  metadata:
    name: svcclusterb
  spec:
    cidr: 10.124.0.0/20
    natOutgoing: false
    disabled: true
EOF

cat > podclusterb.yaml <<EOF
  apiVersion: projectcalico.org/v3
  kind: IPPool
  metadata:
    name: podclusterb
  spec:
    cidr: 10.120.0.0/14
    natOutgoing: false
    disabled: true
EOF
DATASTORE_TYPE=kubernetes KUBECONFIG=cluster-a.yml kubectl calico create -f svcclusterb.yaml
DATASTORE_TYPE=kubernetes KUBECONFIG=cluster-a.yml kubectl calico create -f podclusterb.yaml


cat > allow-tcp-in-cluster-a.yml <<EOF
 apiVersion: projectcalico.org/v3
 kind: GlobalNetworkPolicy
 metadata:
   name: allow-tcp-in-cluster-a
 spec:
   order: 10
   ingress:
     - action: Allow
       protocol: TCP
EOF

cat > allow-tcp-out-cluster-a.yml <<EOF
 apiVersion: projectcalico.org/v3
 kind: GlobalNetworkPolicy
 metadata:
   name: allow-tcp-out-cluster-a
 spec:
   order: 10
   egress:
     - action: Allow
       protocol: TCP
EOF

DATASTORE_TYPE=kubernetes KUBECONFIG=cluster-a.yml kubectl calico apply -f allow-tcp-in-cluster-a.yml
DATASTORE_TYPE=kubernetes KUBECONFIG=cluster-a.yml kubectl calico create -f allow-tcp-out-cluster-a.yml

- Cluster b

gcloud container clusters get-credentials cluster-b --zone=europe-west3-a --project=kmcm-83960

cat > svcclustera.yaml <<EOF
  apiVersion: projectcalico.org/v3
  kind: IPPool
  metadata:
    name: svcclustera
  spec:
    cidr: 10.92.0.0/20
    natOutgoing: false
    disabled: true
EOF

cat > podclustera.yaml <<EOF
  apiVersion: projectcalico.org/v3
  kind: IPPool
  metadata:
    name: podclustera
  spec:
    cidr: 10.88.0.0/14
    natOutgoing: false
    disabled: true
EOF

DATASTORE_TYPE=kubernetes KUBECONFIG=cluster-b.yml kubectl calico create -f svcclustera.yaml
DATASTORE_TYPE=kubernetes KUBECONFIG=cluster-b.yml kubectl calico create -f podclustera.yaml



## Diag Calico

DATASTORE_TYPE=kubernetes KUBECONFIG=cluster-a.yml kubectl calico get node
DATASTORE_TYPE=kubernetes KUBECONFIG=cluster-a.yml kubectl calico get ipPool
DATASTORE_TYPE=kubernetes KUBECONFIG=cluster-b.yml kubectl calico get ipPool
DATASTORE_TYPE=kubernetes KUBECONFIG=cluster-b.yml kubectl calico get nodes


{
    "name": "cluater-a-host-local",
    "cniVersion": "0.1.0",
    "type": "calico",
    "kubernetes": {
        "kubeconfig": "/path/to/kubeconfig",
        "node_name": "node-name-in-k8s"
    },
    "ipam": {
        "type": "host-local",
        "ranges": [
            [
                { "subnet": "usePodCidr" }
            ],
        ],
        "routes": [
            { "dst": "0.0.0.0/0" },
        ]
    }
}
