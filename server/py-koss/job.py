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

cfg = EnvVarConfig()
HTTP_REQUEST_SCHEME = cfg.get("HTTP_REQUEST_SCHEME", str)
SERVICE_IMPORT_ENDPOINT = cfg.get("SERVICE_IMPORT_ENDPOINT", str)
METRICS = cfg.get("METRICS", list)
METRICS_QUERY_PATH = "/api/v1/query"


logger = get_logger()

# ServiceImports Fetching


def import_services_from(json):
    return list(map(lambda item: ServiceImport(item), json.get("items", [])))


def fetch_all_imported_services():
    service_import_api = "http://{}".format(SERVICE_IMPORT_ENDPOINT)
    response = requests.get(service_import_api)
    response.raise_for_status()

    json = response.json()
    return import_services_from(json)

# ClusterMetrics Fetching


def source_to_target_latency_in_seconds(json):
    results = json.get("data", {}).get("result", None)
    if results == None or len(results) == 0:
        raise Exception("source_to_target_latency_in_seconds 1")
    return list(map(lambda res: ClusterMetrics(res), results))


def fetch_metrics_form_all(monitoring_endpoints):
    metrics = []
    for ep in monitoring_endpoints:
        for metric in METRICS:
            api_request_ep = "http://{}/{}".format(ep, METRICS_QUERY_PATH)
            response = requests.get(api_request_ep, params={"query": metric})
            response.raise_for_status()

            json = response.json()
            metrics.append(source_to_target_latency_in_seconds(json))
    return metrics


if __name__ == '__main__':
    # Run TODO!
    imported_services = fetch_all_imported_services()

    monitoring_services = list(
        filter(lambda svc: "prometheus" in svc.name.lower(), imported_services))

    monitoring_endpoints = list(
        map(lambda svc: svc.clsuterset_end_point, monitoring_services))

    metrics = fetch_metrics_form_all(monitoring_endpoints)

    # TODO - build our clusters simiilarly to pass to the optimization and get weights

    logger.info("monitoring_services", monitoring_services)
    logger.info("imported_services", imported_services)
    logger.info("metrics", metrics)
