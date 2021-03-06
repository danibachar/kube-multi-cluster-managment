# TODO
# 1) Fetch all exported/imported prometheus services (this will also indicate the clusters!)

# 2) Fetch data: (https://submariner.io/operations/monitoring/)
# 2.1) demand/between cluster to services - can be fetched from prometheus with query - `submariner_service_discovery_query_counter`
# 2.2) capacity (of each service instance) - TBD! shoud be within the service info data, use annotations
# 2.3) latency between services/clusters- from prometheus using query - `submariner_connection_latency_seconds`

# 3) Understand one hop dependency (might be supply be the user)

# 4) Optional - get payload size of avg request


# 5) Build Models and calculate weights


# 6) Find a way to distribute the weights - maybe using syncers of admiral

from models import Cluster, ServiceImport, ClusterMetrics
from utils import EnvVarConfig, get_logger
import requests

logger = get_logger()

cfg = EnvVarConfig()
HTTP_REQUEST_SCHEME = cfg.get("HTTP_REQUEST_SCHEME", str)
SERVICE_IMPORT_ENDPOINT = cfg.get("SERVICE_IMPORT_ENDPOINT", str)
METRICS = cfg.get("METRICS", list)
METRICS_QUERY_PATH = "/api/v1/query"


def fetch(url, json_parser, params={}):
    logger.info("fetching: {} with params: {}".format(url, params))
    response = requests.get(url, params=params)
    response.raise_for_status()

    json = response.json()
    return json_parser(json)

# ServiceImports Fetching
def fetch_all_imported_services():
    def parser(json):
        return list(map(lambda item: ServiceImport(item), json.get("items", [])))
    service_import_api = "http://{}".format(SERVICE_IMPORT_ENDPOINT)
    return fetch(service_import_api, parser)

# ClusterMetrics Fetching
def fetch_metrics_form_all(monitoring_endpoints):
    def parser(json):
        results = json.get("data", {}).get("result", [])
        return list(map(lambda res: ClusterMetrics(res), results))
    metrics = []
    for ep in monitoring_endpoints:
        for metric in METRICS:
            api_request_ep = "http://{}/{}".format(ep, METRICS_QUERY_PATH)
            metrics.append(
                fetch(api_request_ep, parser, {"query": metric}))
    return metrics


if __name__ == '__main__':
    # Run TODO!
    imported_services = fetch_all_imported_services()

    # For getting metrics and relevant data from all clusters
    monitoring_services = list(
        filter(lambda svc: "prometheus" in svc.name.lower(), imported_services))
    monitoring_endpoints = list(set(
        map(lambda svc: svc.clsuterset_end_point, monitoring_services)))
    # The relevant metrics
    metrics = fetch_metrics_form_all(monitoring_endpoints)
    # For exporting services with new weights per cluster
    exporting_services = list(
        filter(lambda svc: "serviceexporter" in svc.name.lower(), imported_services))
    exporting_endpoints = list(set(
        map(lambda svc: svc.clsuterset_end_point, exporting_services)))

    # TODO - build our clusters simiilarly to pass to the optimization and get weights
    # TODO - after building clusters and running optimization configure new service export with annotations for each service per cluster
    logger.info("exporting_services", exporting_services)
    logger.info("monitoring_services", monitoring_services)
    logger.info("imported_services", imported_services)
    logger.info("metrics", metrics)
    # TODO -
