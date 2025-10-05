#!/usr/bin/env bash
set -euo pipefail
URL_OK=${URL_OK:-http://localhost:8080/}
URL_ERR=${URL_ERR:-http://localhost:8080/error}
N=${N:-50}
E=${E:-15}

while true; do
  echo "Send $N OK requests to $URL_OK and $E ERR requests to $URL_ERR ..."
  for i in $(seq 1 $N); do curl -s "$URL_OK" >/dev/null; done
  for i in $(seq 1 $E); do curl -s "$URL_ERR" >/dev/null; done
  echo "Done. Check Grafana/Prometheus/Alertmanager."
  sleep 0.5
done
