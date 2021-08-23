package utils

import (
	"github.com/danibachar/kube-multi-cluster-managment/server/golang-pkg/models"
)

func GCPSetCurrentWorkingCluster(config models.GCPClusterConfig) (error, string) {
	command := "gcloud container clusters get-credentials " + config.ClusterName
	command = command + " --region=" + config.Region
	command = command + " --project=" + config.Project
	return RunShell(command)
}

// Cluster CIDR is the range of the pods IPs, up to 1008 pods on GCP
func clusterCIDRCommand() string {
	command := "kubectl cluster-info dump | grep -m 1 cluster-cidr | grep -E -o '(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\/(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)'"
	return command
}

func GCPGetClusterCIDR() (error, string) {
	command := clusterCIDRCommand()
	err, res := RunShell(command)
	if err != nil {
		return err, ""
	}
	return nil, res
}

// Service CIDR is the services IP range
func serviceCIDRCommand() string {
	command := `echo '{"apiVersion":"v1","kind":"Service","metadata":{"name":"tst"},"spec":{"clusterIP":"1.1.1.1","ports":[{"port":443}]}}' | kubectl apply -f - 2>&1 | sed 's/.*valid IPs is //'`
	return command
}

func GCPGetServiceCIDR() (error, string) {
	command := serviceCIDRCommand()
	err, res := RunShell(command)
	if err != nil {
		return err, ""
	}
	return nil, res
}
