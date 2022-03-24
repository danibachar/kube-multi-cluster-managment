import http from 'k6/http';
import { sleep } from 'k6';

const ENDPOINT = 'http://34.133.189.105/load'
let headers = { 'Content-Type': 'application/json' };
const PAYLOAD = `{"memory_params": {"duration_seconds": 0.2, "kb_count": 500}, "cpu_params": {"duration_seconds": 0.2, "load": 0.2}}`

export const options = {
    scenarios: {
    //   constant_request_rate_01: {
    //     executor: 'constant-arrival-rate',
    //     rate: 1,
    //     timeUnit: '1s',
    //     duration: '2m',
    //     preAllocatedVUs: 20,
    //     maxVUs: 10000,
    //     startTime: "0s",
    //     gracefulStop: "0s",
    //   },
    //   constant_request_rate_02: {
    //     executor: 'constant-arrival-rate',
    //     rate: 2,
    //     timeUnit: '1s',
    //     duration: '2m',
    //     preAllocatedVUs: 20,
    //     maxVUs: 10000,
    //     startTime: "2m",
    //     gracefulStop: "0s",
    //   },
    //   constant_request_rate_03: {
    //     executor: 'constant-arrival-rate',
    //     rate: 3,
    //     timeUnit: '1s',
    //     duration: '2m',
    //     preAllocatedVUs: 30,
    //     maxVUs: 10000,
    //     startTime: "4m",
    //     gracefulStop: "0s",
    //   },
    //   constant_request_rate_04: {
    //     executor: 'constant-arrival-rate',
    //     rate: 4,
    //     timeUnit: '1s',
    //     duration: '2m',
    //     preAllocatedVUs: 30,
    //     maxVUs: 10000,
    //     startTime: "6m",
    //     gracefulStop: "0s",
    //   },
      constant_request_rate_05: {
        executor: 'constant-arrival-rate',
        rate: 100,
        timeUnit: '1s',
        duration: '5m',
        preAllocatedVUs: 500,
        maxVUs: 10000,
        startTime: "10m",
        gracefulStop: "0s",
      },
    //   constant_request_rate_06: {
    //     executor: 'constant-arrival-rate',
    //     rate: 6,
    //     timeUnit: '1s',
    //     duration: '2m',
    //     preAllocatedVUs: 40,
    //     maxVUs: 10000,
    //     startTime: "10m",
    //     gracefulStop: "0s",
    //   },
    //   constant_request_rate_07: {
    //     executor: 'constant-arrival-rate',
    //     rate: 7,
    //     timeUnit: '1s',
    //     duration: '2m',
    //     preAllocatedVUs: 50,
    //     maxVUs: 10000,
    //     startTime: "12m",
    //     gracefulStop: "0s",
    //   },
    //   constant_request_rate_08: {
    //     executor: 'constant-arrival-rate',
    //     rate: 8,
    //     timeUnit: '1s',
    //     duration: '2m',
    //     preAllocatedVUs: 50,
    //     maxVUs: 10000,
    //     startTime: "14m",
    //     gracefulStop: "0s",
    //   },
    //   constant_request_rate_09: {
    //     executor: 'constant-arrival-rate',
    //     rate: 9,
    //     timeUnit: '1s',
    //     duration: '2m',
    //     preAllocatedVUs: 60,
    //     maxVUs: 10000,
    //     startTime: "16m",
    //     gracefulStop: "0s",
    //   },
    //   constant_request_rate_010: {
    //     executor: 'constant-arrival-rate',
    //     rate: 10,
    //     timeUnit: '1s',
    //     duration: '2m',
    //     preAllocatedVUs: 60,
    //     maxVUs: 10000,
    //     startTime: "18m",
    //     gracefulStop: "0s",
    //   },
    },
};

// export const options = {
//     stages: [
//         { rate: 1, timeUint: '1s', duration: '5m', target: 10 },
//         { rate: 2, timeUint: '1s', duration: '5m', target: 20 },
//         { rate: 3, timeUint: '1s', duration: '5m', target: 30 },
//     ],
// };

export default function () {
  http.post(ENDPOINT, PAYLOAD, { headers: headers }  );
}

