import http from 'k6/http';
import { sleep } from 'k6';

const ENDPOINT = 'http://34.133.189.105/load'
let headers = { 'Content-Type': 'application/json' };
const PAYLOAD = `{"memory_params": {"duration_seconds": 0.2, "kb_count": 500}, "cpu_params": {"duration_seconds": 0.2, "load": 0.2}}`

export const options = {
    scenarios: {
      constant_request_rate_05: {
        executor: 'constant-arrival-rate',
        rate: 200,
        timeUnit: '1s',
        duration: '5m',
        preAllocatedVUs: 1000,
        maxVUs: 10000,
        startTime: "10m",
        gracefulStop: "0s",
      },
    },
};

export default function () {
  http.post(ENDPOINT, PAYLOAD, { headers: headers }  );
}

