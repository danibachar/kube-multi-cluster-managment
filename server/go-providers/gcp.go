package main

import (
	"strings"

	"github.com/danibachar/kube-multi-cluster-managment/server/golang-pkg/models"
	"github.com/danibachar/kube-multi-cluster-managment/server/golang-pkg/utils"
	"github.com/labstack/gommon/log"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

func clusterCreateCommand(config models.GCPClusterConfig) string {
	command := "gcloud container clusters create " + config.ClusterName
	command = command + " --zone=" + config.Region
	command = command + " --project=" + config.Project
	command = command + " --machine-type=" + config.NodeMachineType
	command = command + " --num-nodes=" + config.WorkerMachineCount
	// command = command+"--node-version="+config.KubeVersion
	// command = command + " --cluster-version=" + config.KubeVersion

	return command
}

func createCluster(config models.GCPClusterConfig) (error, string) {
	log.Info("start cluster creation")
	command := clusterCreateCommand(config)
	return utils.RunShell(command)
}

// Submariner workaround scripts
func workaroundRPFiltersCommand(config models.GCPClusterConfig) string {
	command := "./configure-rp-filter.sh"
	return command
}

func applySubmarinerGCPWorkAround(config models.GCPClusterConfig) error {
	log.Info("applying gcp submariner networking workaround")
	command := workaroundRPFiltersCommand(config)
	err, res := utils.RunShell(command)
	log.Info("gcp submariner networking workaround finished with %v", res)
	return err
}

// Firewall
func firewallTCPRuleCommand(commandType, name, direction, clusterCIDR, serviceCIDR, project string) string {
	var rangeFlag string = "--source-ranges="
	if direction == "OUT" {
		rangeFlag = "--destination-ranges="
	}
	components := []string{"gcloud compute firewall-rules ", "--project=", project, " ", commandType, " ", name, " --allow=tcp ", rangeFlag, clusterCIDR, ",", serviceCIDR}
	if commandType == "create" {
		components = append(components, " --direction=")
		components = append(components, direction)
	}
	command := strings.Join(components, "")
	return command
}

func firewallUDPRuleCommand(name, port, direction, project string) string {
	components := []string{"gcloud compute firewall-rules create ", "--project=", project, " ", name, " --allow=udp:", port, " ", "--direction=", direction}
	command := strings.Join(components, "")
	return command
}

func isFirewallRuleExists(name, project string) bool {
	command := "gcloud compute firewall-rules list"
	command = command + " --project=" + project
	command = command + " | grep " + name
	err, res := utils.RunShell(command)
	if err != nil {
		return false
	}
	return res != ""
}

func createOrUpdateFirewallRules(config models.GCPClusterConfig) error {
	err, clusterCIDR := utils.GCPGetClusterCIDR()
	if err != nil {
		return err
	}
	log.Info("clusterCIDR ", clusterCIDR)

	// time.Sleep(1000 * time.Millisecond)
	err, serviceCIDR := utils.GCPGetServiceCIDR()
	if err != nil {
		return err
	}
	log.Info("serviceCIDR ", serviceCIDR)

	// time.Sleep(1000 * time.Millisecond)
	inCommandName := "allow-tcp-in-" + config.ClusterName + "-" + config.Region
	var inCommandType string
	if isFirewallRuleExists(inCommandName, config.Project) {
		inCommandType = "update"
	} else {
		inCommandType = "create"
	}

	// time.Sleep(1000 * time.Millisecond)
	command := firewallTCPRuleCommand(inCommandType, inCommandName, "IN", clusterCIDR, serviceCIDR, config.Project)
	err, res := utils.RunShell(command)
	if err != nil {
		return err
	}
	log.Info("% finished with", inCommandName, res)

	// time.Sleep(1000 * time.Millisecond)
	outCommandName := "allow-tcp-out-" + config.ClusterName + "-" + config.Region
	var outCommandType string
	if isFirewallRuleExists(outCommandName, config.Project) {
		outCommandType = "update"
	} else {
		outCommandType = "create"
	}

	// time.Sleep(1000 * time.Millisecond)
	command = firewallTCPRuleCommand(outCommandType, outCommandName, "OUT", clusterCIDR, serviceCIDR, config.Project)
	err, res = utils.RunShell(command)
	if err != nil {
		return err
	}
	log.Info("% finished with", outCommandName, res)
	return nil
}

// Public
func GCPCreateOrUpdateCluster(config models.GCPClusterConfig) error {
	// TODO - update or create - check if cluster exists
	err, res := createCluster(config)
	if err != nil {
		return err
	}
	log.Info("cluster created with", res)

	err, res = utils.GCPSetCurrentWorkingCluster(config)
	if err != nil {
		return err
	}
	log.Info("cluster %s was configured with %v", config.ClusterName, res)

	if err := applySubmarinerGCPWorkAround(config); err != nil {
		return err
	}

	if err := createOrUpdateFirewallRules(config); err != nil {
		return err
	}

	return nil
}

func deleteCluster(config models.GCPClusterConfig) (error, string) {
	command := "gcloud container clusters delete " + config.ClusterName
	command = command + " --region=" + config.Region
	command = command + " --project=" + config.Project
	command = command + " --quiet"
	return utils.RunShell(command)
}

func GCPDeleteCluster(config models.GCPClusterConfig) error {
	log.Info("start cluster deletion %s", config.ClusterName)
	err, _ := deleteCluster(config)
	if err != nil {
		return err
	}
	log.Info("cluster %s deleted succesfuly", config.ClusterName)
	return nil
}

func GetGCPRegions(project string) (error, []*compute.Region) {
	ctx := context.Background()

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
		return err, nil
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
		return err, nil
	}

	regions := make([]*compute.Region, 0)
	req := computeService.Regions.List(project)
	if err := req.Pages(ctx, func(page *compute.RegionList) error {
		newRegions := page.Items
		regions = append(regions, newRegions...)
		// for _, region := range page.Items {

		// }
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	return nil, regions
}
