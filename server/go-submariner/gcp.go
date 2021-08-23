package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/danibachar/kube-multi-cluster-managment/server/providers/pkg/models"
	"github.com/danibachar/kube-multi-cluster-managment/server/providers/pkg/utils"
	"github.com/labstack/gommon/log"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

const (
	baseSubmarinerCLI      string = "~/.local/bin/subctl"
	submarinerGatewayLabel        = "submariner.io/gateway"
	trueLabel                     = "true"
)

func GCPDeploySubmarinerBrokerOn(config models.GCPClusterConfig) error {
	err, res := utils.GCPSetCurrentWorkingCluster(config)
	if err != nil {
		return err
	}
	log.Info("set cluster with res", res)
	command := baseSubmarinerCLI + " deploy-broker"
	err, res = utils.RunBash(command)
	if err != nil {
		return err
	}
	// TODO - we need to handle the broker-info.subm, write it in the DB so we can handle the broker cluster
	log.Info("Deployed broker on clusrer ", config.ClusterName)
	return nil
}

func GCPJoinClusterToBroker(config models.GCPClusterConfig) error {
	err, res := utils.GCPSetCurrentWorkingCluster(config)
	if err != nil {
		return err
	}
	log.Info("set cluster with res", res)

	err = addNeededLabelToFirstAvailableNode()
	if err != nil {
		return err
	}
	log.Info("enabled gateway labels for node")

	err, serviceCIDR := utils.GCPGetServiceCIDR()
	if err != nil {
		return err
	}
	command := baseSubmarinerCLI + " join broker-info.subm"
	command = command + " --clusterid " + config.ClusterName
	command = command + " --servicecidr " + serviceCIDR
	err, res = utils.RunBash(command)
	if err != nil {
		return err
	}

	log.Info("Cluster %s Joined cluster set", config.ClusterName)
	return nil
}

func addNeededLabelToFirstAvailableNode() error {
	clientset, err := utils.GetKubernetesClient()
	if err != nil {
		return err
	}
	log.Info("clusterset is configured")
	var node string
	node, err = getNodeForLabeling(clientset)
	if err != nil {
		return err
	}
	if node == "" { // hack should return error and create an error for already labeled node
		log.Info("Some node already labeled")
		return nil
	}
	log.Info("labeling node %s", node)
	return addLabelsToNode(clientset, node, map[string]string{submarinerGatewayLabel: trueLabel})
}

func getNodeForLabeling(clientset kubernetes.Interface) (string, error) {
	// List Submariner-labeled nodes
	selector := labels.SelectorFromSet(labels.Set(map[string]string{submarinerGatewayLabel: trueLabel}))
	labeledNodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return "", err
	}
	if len(labeledNodes.Items) > 0 {
		fmt.Printf("DB: * There are %d labeled nodes in the cluster:\n", len(labeledNodes.Items))
		for _, node := range labeledNodes.Items {
			fmt.Printf("  - %s\n", node.GetName())
		}
		return "", nil
	}
	// List the worker nodes and select one
	workerNodes, err := clientset.CoreV1().Nodes().List(
		context.TODO(), metav1.ListOptions{LabelSelector: "node-role.kubernetes.io/worker"})
	if err != nil {
		return "", err
	}
	if len(workerNodes.Items) == 0 {
		// In some deployments (like KIND), worker nodes are not explicitly labelled. So list non-master nodes.
		workerNodes, err = clientset.CoreV1().Nodes().List(
			context.TODO(), metav1.ListOptions{LabelSelector: "!node-role.kubernetes.io/master"})
		if err != nil {
			return "", err
		}
		if len(workerNodes.Items) == 0 {
			// No worker
			return "", nil
		}
	}
	// Chossing the first node in the list, for no particular reason
	return workerNodes.Items[0].GetName(), nil
}

// this function was sourced from:
// https://github.com/kubernetes/kubernetes/blob/a3ccea9d8743f2ff82e41b6c2af6dc2c41dc7b10/test/utils/density_utils.go#L36
func addLabelsToNode(c kubernetes.Interface, nodeName string, labelsToAdd map[string]string) error {
	var tokens = make([]string, 0, len(labelsToAdd))
	for k, v := range labelsToAdd {
		tokens = append(tokens, fmt.Sprintf("\"%s\":\"%s\"", k, v))
	}

	labelString := "{" + strings.Join(tokens, ",") + "}"
	patch := fmt.Sprintf(`{"metadata":{"labels":%v}}`, labelString)

	// retry is necessary because nodes get updated every 10 seconds, and a patch can happen
	// in the middle of an update

	var lastErr error
	err := wait.ExponentialBackoff(nodeLabelBackoff, func() (bool, error) {
		_, lastErr = c.CoreV1().Nodes().Patch(context.TODO(), nodeName, types.MergePatchType, []byte(patch), metav1.PatchOptions{})
		if lastErr != nil {
			if !errors.IsConflict(lastErr) {
				return false, lastErr
			}
			return false, nil
		} else {
			return true, nil
		}
	})

	if err == wait.ErrWaitTimeout {
		return lastErr
	}

	return err
}

var nodeLabelBackoff wait.Backoff = wait.Backoff{
	Steps:    10,
	Duration: 1 * time.Second,
	Factor:   1.2,
	Jitter:   1,
}
