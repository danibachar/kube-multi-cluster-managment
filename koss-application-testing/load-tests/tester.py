#!/usr/bin/env python3

from time import sleep
import pandas as pd
import numpy as np
import requests

DEBUG = False
def _query(query):
  if DEBUG:
    print("query {}".format(query))
  url = "http://localhost:8082/api/v1/query?query={}".format(query)
  response = requests.get(url).json()
  if DEBUG:
    print("response {}".format(response))
  return _extract_result(response=response)

def _extract_result(response):
  # {"status":"success","data":{"resultType":"vector","result":[{"metric":{"instance":"simple-svc:80","job":"simple-svc-scraping","method":"POST","path":"/load","status":"200"},"value":[1647843667.065,"0.4666977798519901"]}]}}
  result = response.get("data",{}).get("result",[])
  if len(result) > 0:
    value = result[0].get("value",[])
    if len(value) == 2:
      return float(value[1])
    else:
      if DEBUG:
        print("query result has empty value")  
  else:
    if DEBUG:
      print("query returned empty result")
  return np.nan

def _validate_extraction(value):
  return (value == 0 or value == 0.0 or value == np.nan or value == None)

def extract_rps():
  query = "rate(flask_http_request_duration_seconds_count{status='200'}[30s])"
  return _query(query=query)

def extract_mean_rtt():
  query = "rate(flask_http_request_duration_seconds_sum{status='200'}[30s])/rate(flask_http_request_duration_seconds_count{status='200'}[30s])"
  return _query(query=query)

def extract_90_percentile_rtt():
  query = "histogram_quantile(0.9, rate(flask_http_request_duration_seconds_bucket{status='200'}[30s]))"
  return _query(query=query)

def extract_95_percentile_rtt():
  query = "histogram_quantile(0.95, rate(flask_http_request_duration_seconds_bucket{status='200'}[30s]))"
  return _query(query=query)

def extract_99_percentile_rtt():
  query = "histogram_quantile(0.99, rate(flask_http_request_duration_seconds_bucket{status='200'}[30s]))"
  return _query(query=query)

def extract_errors_per_second():
  query = "sum(rate(flask_http_request_duration_seconds_count{status!='200'}[30s]))"
  return _query(query=query)


if __name__ == '__main__':
  rpses = []
  mean_rttes = []
  per_90_rttes = []
  per_95_rttes = []
  per_99_rttes = []
  errors_per_second = []
  while True:
    rps = extract_rps()
    if _validate_extraction(rps):
      print("skipping empty rps")
      continue
    mean_rtt = extract_mean_rtt()
    if _validate_extraction(mean_rtt):
      print("skipping empty mean_rtt")
      continue
    per_90_rtt = extract_90_percentile_rtt()
    if _validate_extraction(per_90_rtt):
      print("skipping empty per_90_rtt")
      continue
    per_95_rtt = extract_95_percentile_rtt()
    if _validate_extraction(per_95_rtt):
      print("skipping empty per_95_rtt")
      continue
    per_99_rtt = extract_99_percentile_rtt()
    if _validate_extraction(per_99_rtt):
      print("skipping empty per_99_rtt")
      continue
    errors_per_second.append(extract_errors_per_second())
    rpses.append(rps)
    mean_rttes.append(mean_rtt)
    per_90_rttes.append(per_90_rtt)
    per_95_rttes.append(per_95_rtt)
    per_99_rttes.append(per_99_rtt)
    d = {
      'rps': rpses, 
      'mean_rtt': mean_rtt,
      'per_90_rtt': per_90_rtt, 
      'per_95_rtt': per_95_rtt, 
      'per_99_rtt': per_99_rtt,
      "errors_per_second": errors_per_second
    }
    df = pd.DataFrame(data=d)
    df.to_csv('rps-to-rtt.csv', index=False)
    sleep(0.5)

