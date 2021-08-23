from timeit import default_timer as timer
import numpy as np
import math
from pulp import *


def simple_max_addative_weight(
    price,
    max_price,
    price_weight,
    latency,
    max_latency,
    latency_weight
):
    if max_price == 0:
        raise
    if max_latency == 0:
        raise

    price_part = (price / max_price) * price_weight
    latency_part = (latency / max_latency) * latency_weight

    return price_part + latency_part

############################################################


def _services(clusters, target_type):
    source_services = []
    destination_services = []
    for cluster in clusters:
        for service in cluster.services.values():
            if target_type in [d.job_type for d in service.dependencies]:
                source_services.append(service)
            if target_type == service.job_type:
                destination_services.append(service)
    return source_services, destination_services


def _lp_params(source_services, destination_services, cost_function_weights):
    x = []
    for src in source_services:
        x.append([])
        my_propogated_request_count = math.floor(
            src.jobs_consumed_per_time_slot)
        for dest in destination_services:
            name = src.full_name+"-"+dest.full_name
            # Continuous # Integer
            var = pulp.LpVariable(
                name, lowBound=0, upBound=my_propogated_request_count, cat='Continuous')
            x[-1].append(var)

    capacities = []
    for dest in destination_services:
        capacities.append(dest.capacity)

    demands = []
    for src in source_services:
        my_propogated_request_count = math.floor(
            src.jobs_consumed_per_time_slot)
        demands.append(my_propogated_request_count)

    costs = []
    price_costs = []
    latency_costs = []
    for src in source_services:
        prices = []
        latencies = []
        for dest in destination_services:
            prices.append(src.cluster.zone.price_per_gb(dest.cluster.zone))
            latencies.append(
                src.cluster.zone.latency_per_request(dest.cluster.zone))
        max_price = max(prices)
        max_latency = max(latencies)

        min_price = min(prices)
        min_latency = min(latencies)
        dests_cost = [simple_min_addative_weight(
            price=p,
            min_price=min_price,
            price_weight=cost_function_weights[0],
            latency=l,
            min_latency=min_latency,
            latency_weight=cost_function_weights[1]
        ) for p, l in zip(prices, latencies)]

        price_costs.append(np.asarray(prices) / max_price)
        latency_costs.append(np.asarray(latencies) / max_latency)

        costs.append(dests_cost)

    return x, costs, demands, capacities, price_costs, latency_costs


def _build_problem(source_services, destination_services, x, costs, demands, capacities, price_costs, latency_costs, cost_function_weights):
    # Define Optimization Problem - Minimizing Cost
    prob = LpProblem("Service Selection Problem", LpMinimize)
    src_count = len(source_services)
    dest_count = len(destination_services)
    # Objective
    prob += lpSum(x[src_idx][dst_idx] * costs[src_idx][dst_idx]
                  for src_idx in range(src_count) for dst_idx in range(dest_count))
    # prob += lpSum(x[src_idx][dst_idx]*cost_function_weights[0]*price_costs[src_idx][dst_idx] + x[src_idx][dst_idx]*cost_function_weights[1]*latency_costs[src_idx][dst_idx] for src_idx in range(src_count) for dst_idx in range(dest_count))
    # Constriants:
    # (1) All Request demand must be deplited - i.e all requests must be sent and cannot be dropped
    for src_idx in range(src_count):
        prob += lpSum([x[src_idx][dst_idx]
                      for dst_idx in range(dest_count)]) == demands[src_idx]
    # (2) Capacity - service cannot handle more than its capacity
    for dst_idx in range(dest_count):
        prob += lpSum([x[src_idx][dst_idx]
                      for src_idx in range(src_count)]) <= capacities[dst_idx]
        # (3) # Liveness
    # for dst_srv_idx in dest_srvs_idx:
    #     prob += lpSum([x[src_srv_idx][dst_srv_idx]]) >= 1

    prob.writeLP("ServiceSelection.lp")
    return prob


def _weights_distribution(x, source_services, destination_services):
    res = {}
    src_count = len(source_services)
    dest_count = len(destination_services)

    for src_svc_idx, src_svc in enumerate(source_services):
        if src_svc.cluster.id not in res:
            res[src_svc.cluster.id] = {}
        if src_svc.id not in res[src_svc.cluster.id]:
            res[src_svc.cluster.id][src_svc.id] = {}
        dists = [x[src_svc_idx][dst_svc_idx].value()
                 for dst_svc_idx in range(dest_count)]
        estimated_total_requests = sum(dists)
        if estimated_total_requests == 0:
            print("estimated_total_requests were zero, pading with 1")
            estimated_total_requests = 1
        for dst_svc_idx in range(dest_count):
            dst_svc = destination_services[dst_svc_idx]
            if dst_svc.job_type not in res[src_svc.cluster.id][src_svc.id]:
                res[src_svc.cluster.id][src_svc.id][dst_svc.job_type] = {}
            if dst_svc.cluster in res[src_svc.cluster.id][src_svc.id][dst_svc.job_type]:
                # print(res[src_svc.cluster.id][dst_svc.job_type])
                # print(src_svc.cluster.id)
                # print(dst_svc.job_type)
                # print(dst_svc.cluster)
                raise
            res[src_svc.cluster.id][src_svc.id][dst_svc.job_type][dst_svc.cluster] = x[src_svc_idx][dst_svc_idx].value(
            ) / estimated_total_requests

    return res


def calculate_weights_for(clusters, cost_function_weights):
    res = {}  # Map of maps, cluster.id <-> weights map
    target_types = list(sorted(set([dependency.job_type for cluster in clusters for service in cluster.services.values(
    ) for dependency in service.dependencies])))
    # Optimize per target type
    for target_type in target_types:
        sources, destinations = _services(clusters, target_type)

        if len(destinations) == 0:
            print("target_type", target_type)
            print("sources", sources)
            print("destinations", destinations)
            continue
        x, costs, demands, capacities, price_costs, latency_costs = _lp_params(
            sources, destinations, cost_function_weights)
        prob = _build_problem(sources, destinations, x, costs, demands,
                              capacities, price_costs, latency_costs, cost_function_weights)
        s = timer()
        prob.solve(PULP_CBC_CMD(msg=0))
        e = timer()
        # print("solving took {} seconds".format(e-s))
        w = _weights_distribution(x, sources, destinations)
        # print("w for {} = {}".format(target_type, w))
        # res = {**res, **w}
        res[target_type] = w
        # print("res after update = {}".format(res))
        if prob.status == -1:
            print("###########################")
            print("could not solve {}, at tik = {}".format(target_type, at_tik))
            print("demands = {}".format(demands))
            print("capacities = {}".format(capacities))
            print("###########################")
    # print(res)
    return res
