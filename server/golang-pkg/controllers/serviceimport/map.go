/*
SPDX-License-Identifier: Apache-2.0
Copyright Contributors to the Submariner project.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package serviceimport

import (
	"fmt"
	"sync"

	lhconstants "github.com/submariner-io/lighthouse/pkg/constants"
	mcsv1a1 "sigs.k8s.io/mcs-api/pkg/apis/v1alpha1"
)

type ServiceInfo struct {
	IP          string
	Ports       []mcsv1a1.ServicePort
	HostName    string
	ClusterName string
}

type clusterInfo struct {
	name     string
	services map[string]*ServiceInfo
}

type namespaceInfo struct {
	name     string
	clusters map[string]*clusterInfo
}

type Map struct {
	svcMap map[string]*namespaceInfo
	sync.RWMutex
}

func (m *Map) getAllExportedServicesIn(namespace string) ([]*ServiceInfo, error) {
	// Checking namespace
	namespaceInfo, namespaceInfoExists := m.svcMap[namespace]
	if !namespaceInfoExists {
		return nil, fmt.Errorf("namespace does not exists")
	}

	services := make([]*ServiceInfo, 0)
	for _, clusterInfo := range namespaceInfo.clusters {
		for _, serviceInfo := range clusterInfo.services {
			services = append(services, serviceInfo)
		}
	}
	return services, nil
}

func (m *Map) GetAllExportedServicesIn(namespace string) ([]*ServiceInfo, error) {
	m.RLock()
	defer m.RUnlock()
	return m.getAllExportedServicesIn(namespace)
}

func (m *Map) GetAllExportedServices() ([]*ServiceInfo, error) {
	m.RLock()
	defer m.RUnlock()

	services := make([]*ServiceInfo, 0)
	for _, namespaceInfo := range m.svcMap {
		namespace := namespaceInfo.name
		if serviceInfos, err := m.getAllExportedServicesIn(namespace); err != nil {
			services = append(services, serviceInfos...)
		}
	}

	return services, nil
}

func NewMap() *Map {
	return &Map{
		svcMap: make(map[string]*namespaceInfo),
	}
}

func (m *Map) Put(serviceImport *mcsv1a1.ServiceImport) {
	if name, ok := serviceImport.Annotations["origin-name"]; ok {
		namespace := serviceImport.Annotations["origin-namespace"]

		m.Lock()
		defer m.Unlock()

		// Checking namespace
		nsInfo, namespaceInfoExists := m.svcMap[namespace]
		if !namespaceInfoExists {
			nsInfo = &namespaceInfo{
				name:     namespace,
				clusters: make(map[string]*clusterInfo),
			}
		}

		// Checking clsuters
		clusterName := serviceImport.GetLabels()[lhconstants.LabelSourceCluster]
		csInfo, clusterInfoExists := nsInfo.clusters[clusterName]
		if !clusterInfoExists {
			csInfo = &clusterInfo{
				name:     clusterName,
				services: make(map[string]*ServiceInfo),
			}
		}

		if serviceImport.Spec.Type == mcsv1a1.ClusterSetIP {
			sInfo := &ServiceInfo{
				IP:          serviceImport.Spec.IPs[0],
				Ports:       serviceImport.Spec.Ports,
				HostName:    name,
				ClusterName: serviceImport.GetLabels()[lhconstants.LabelSourceCluster],
			}

			csInfo.services[name] = sInfo

		}

		nsInfo.clusters[clusterName] = csInfo
		m.svcMap[namespace] = nsInfo
	}
}

func (m *Map) Remove(serviceImport *mcsv1a1.ServiceImport) {
	if name, ok := serviceImport.Annotations["origin-name"]; ok {
		namespace := serviceImport.Annotations["origin-namespace"]

		m.Lock()
		defer m.Unlock()

		namespaceInfo, namespaceInfoExists := m.svcMap[namespace]
		if !namespaceInfoExists {
			return
		}

		for _, info := range serviceImport.Status.Clusters {
			clusterInfo, clusterInfoExists := namespaceInfo.clusters[info.Cluster]
			if !clusterInfoExists {
				continue
			}
			delete(clusterInfo.services, name)

			if len(clusterInfo.services) == 0 {
				delete(namespaceInfo.clusters, info.Cluster)
			}
		}

		if len(namespaceInfo.clusters) == 0 {
			delete(m.svcMap, namespace)
		}
	}
}
