docker-run-simple-test:
	docker run --rm -i grafana/k6 run --vus 10 --duration 30s - simple-test.js

run-simple-test:
	k6 run ./simple-test.js --summary-export simple.json

run-constant-rps-test:
	k6 run ./constant-rps.js

run-ramp-up:
	k6 run ./ramp-up-tests/01.js --summary-export ./ramp-up-tests/01.json
	k6 run ./ramp-up-tests/05.js --summary-export ./ramp-up-tests/05.json
	k6 run ./ramp-up-tests/05.js --summary-export ./ramp-up-tests/010.json
