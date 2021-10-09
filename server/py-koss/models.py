import collections
import logging


class Cluster:
    """Representing a Kubernetes Cluster"""

    def __init__(self, id):
        self.id = id
        self.services = collections.OrderedDict()

    def __repr__(self):
        return self.zone.__repr__()

    def __str__(self):
        return str(self.__class__) + ": " + str(self.__dict__)

    def __hash__(self):
        h = hash((self.id, self.zone))
        return h

    def __eq__(self, other):
        return self.id == other.id

    def service(self, job_type):
        return self.services.get(job_type, None)

    def add_service(self, service):
        if self.service(service.job_type) != None:
            logging.error(
                "trying to add service with job type that already exists")
            return
        service.add_to_cluster(self)  # settings cluster <-> service connection
        self.services[service.job_type] = service


class ClusterMetrics:
    def __init__(self, dic):
        metric = dic.get("metric", {})
        latency_in_seconds = dic.get("value", [0, float('inf')])[1]
        self.latncy_in_ms = float(latency_in_seconds) * 1000
        self.source = metric.get("local_cluster", None)
        self.remote = metric.get("remote_cluster", None)

    def __str__(self) -> str:
        return self.__repr__

    def __repr__(self):
        return "The latency from {} to {} is {} ms".format(self.source, self.remote, self.latncy_in_ms)


class ServiceImport:
    def __init__(self, dic):
        metadata = dic.get("metadata", None)
        self.cluster_name, self.name, self.namespace = self._service_details_from_metadata(
            metadata)
        self.cluster_ip = self._service_cluster_ip_from(metadata)
        self.port = self._service_port_from(dic)

    def _service_details_from_metadata(self, metadata):
        labels = metadata.get("labels", None)
        if labels == None:
            raise Exception("service_details_from_metadata")
        return labels.get("lighthouse.submariner.io/sourceCluster", None), labels.get("lighthouse.submariner.io/sourceName", None), labels.get("lighthouse.submariner.io/sourceNamespace", None)

    def _service_port_from(self, item):
        spec = item.get("spec", None)
        if spec == None:
            raise Exception("service_port_from 1")
        ports = spec.get("ports", None)
        if ports == None or len(ports) == 0:
            raise Exception("service_port_from 2")
        return ports[0].get("port", None)

    def _service_cluster_ip_from(self, metadata):
        annotations = metadata.get("annotations", None)
        if annotations == None:
            raise Exception("service_cluster_ip_from")
        return annotations.get("cluster-ip", None)

    @property
    def clsuterset_end_point(self):
        return "{}.{}.{}.svc.clusterset.local:{}".format(self.cluster_name, self.name, self.namespace, self.port)

    def __str__(self) -> str:
        return self.__repr__

    def __repr__(self):
        return self.clsuterset_end_point
