module github.com/danibachar/kube-multi-cluster-managment/server/golang-pkg

go 1.14

require (
	github.com/labstack/gommon v0.3.0
	github.com/submariner-io/lighthouse v0.10.1
	k8s.io/api v0.22.2 // indirect
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v1.0.0
	sigs.k8s.io/mcs-api v0.1.0
)

// Pinned to kubernetes-1.19.10
replace (
	k8s.io/api => k8s.io/api v0.19.10
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.10
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.10
	k8s.io/client-go => k8s.io/client-go v0.19.10
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.19.10
)
