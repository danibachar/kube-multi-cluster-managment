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

type AWSClusterConfig struct {
	ClusterName              string `json:"clusterName"`
	ControlPlaneMachineType  string `json:"controlPlaneMachineType"`
	NodeMachineType          string `json:"nodeMachineType"`
	ControlPlaneMachineCount string `json:"controlPlaneMachineCount"`
	WorkerMachineCount       string `json:"workerMachineCount"`
	Region                   string `json:"region"`
	KubeVersion              string `json:"kubeVersion"`
}

// Example
// {
//     "clusterName": "test",
//     "controlPlaneMachineType": "tg4.micro",
//     "nodeMachineType": "tg4.micro",
//     "controlPlaneMachineCount": "3",
//     "workerMachineCount": "3",
//     "region": "us-east-1",
//     "kubeVersion": "v1.22.0"
// }

func CreateAWSCluster(config AWSClusterConfig) error {

	cmd := exec.Command("bash", "-c",
		"export"+" AWS_REGION="+config.Region+" AWS_SSH_KEY_NAME=default"+" AWS_CONTROL_PLANE_MACHINE_TYPE="+config.ControlPlaneMachineType+" AWS_NODE_MACHINE_TYPE="+config.NodeMachineType,
		" && clusterctl "+"generate " + "cluster "+config.ClusterName
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
