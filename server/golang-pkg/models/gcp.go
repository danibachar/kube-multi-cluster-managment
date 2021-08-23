package models

type NodeAutoscaleConfig struct {
	Enabled      bool `json:"enabled"`
	MaxNodeCount int  `json:"maxNodeCount"`
	MinNodeCount int  `json:"minNodeCount"`
}

type GCPClusterConfig struct {
	ClusterName        string              `json:"clusterName"`
	NodeMachineType    string              `json:"nodeMachineType"`
	WorkerMachineCount string              `json:"workerMachineCount"`
	Region             string              `json:"region"`
	KubeVersion        string              `json:"kubeVersion"`
	Project            string              `json:"project"`
	Autoscale          NodeAutoscaleConfig `json:"autoscale"`
}

type GCPRegion struct {
	Name  string   `json:"name"`
	Zones []string `json:"name"`
}

type GCPRegions struct {
	Regions []GCPRegion `json:"regions"`
}
