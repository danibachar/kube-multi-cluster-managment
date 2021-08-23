package main

import (
	"fmt"
	"os/exec"
)

func bashEcho(value string) error {
	echo := exec.Command("bash", "-c", "echo "+"$"+value)

	stdout, err := echo.Output()
	if err != nil {
		return err
	}
	fmt.Printf("echo stdout", string(stdout))
	return nil
}

func export(key, value string) error {
	exports := exec.Command("bash", "-c", "export "+key+"="+value+"&& echo $"+key)

	stdout, err := exports.Output()
	if err != nil {
		return err
	}
	fmt.Printf("export stdout", string(stdout))
	return bashEcho(key)
}

type GCPClusterConfig struct {
	ClusterName              string `json:"clusterName"`
	ControlPlaneMachineType  string `json:"controlPlaneMachineType"`
	NodeMachineType          string `json:"nodeMachineType"`
	ControlPlaneMachineCount string `json:"controlPlaneMachineCount"`
	WorkerMachineCount       string `json:"workerMachineCount"`
	Region                   string `json:"region"`
	KubeVersion              string `json:"kubeVersion"`
	Project 				 string `json:"project"`
}

// Example
// {
//     "clusterName": "test",
//     "controlPlaneMachineType": "n1-standard-2",
//     "nodeMachineType": "n1-standard-2",
//     "controlPlaneMachineCount": "3",
//     "workerMachineCount": "3",
//     "region": "us-east1",
//     "kubeVersion": "v1.22.0",
//	   "project": "kmcm-83960"
// }

func CreateGCPCluster(config GCPClusterConfig) error {

	cmd := exec.Command("bash", "-c",
		"export"
		+" GCP_REGION="+config.Region
		+" GCP_PROJECT="+config.Project
		+" GCP_NETWORK_NAME=default"
		+" GCP_CONTROL_PLANE_MACHINE_TYPE="+config.ControlPlaneMachineType
		+" GCP_NODE_MACHINE_TYPE="+config.NodeMachineType
		+ " && clusterctl "+"generate " + "cluster "+config.ClusterName
		+" --kubernetes-version "+config.KubeVersion,
		+" --control-plane-machine-count="+config.ControlPlaneMachineCount
		+" --worker-machine-count="+config.WorkerMachineCount
		+" | kubectl apply -f -",
	)
	stdout, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Printf("clusterctl output", string(stdout))
	return nil
}

clusterctl generate cluster test1 --kubernetes-version v1.22.0 --control-plane-machine-count=3 --worker-machine-count=3 | kubectl apply -f -