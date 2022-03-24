import http from 'k6/http';
import { sleep } from 'k6';

let headers = { 'Content-Type': 'application/json' };
const ENDPOINT = 'http://34.133.189.105/load'
const PAYLOAD = `{"memory_params": {"duration_seconds": 0.2, "kb_count": 50}, "cpu_params": {"duration_seconds": 0.2, "load": 0.2}}`

export const options = {
  vus: 10,
  duration: '30s',
};

export default function () {
    const res = http.post(ENDPOINT, PAYLOAD, { headers: headers }  );
    console.log("res", JSON.stringify(res));
    sleep(1);
}