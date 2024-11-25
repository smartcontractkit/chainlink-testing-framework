#!/bin/bash

CONTAINER_NAME="compose-grafana-1"
LOG_MESSAGE="HTTP Server Listen"
MAX_ATTEMPTS=10
attempt=0

# Wait for the log message
until docker logs "$CONTAINER_NAME" 2>&1 | grep -q "$LOG_MESSAGE"; do
    attempt=$((attempt + 1))
    if [ "$attempt" -ge "$MAX_ATTEMPTS" ]; then
        echo "Max attempts reached. Grafana not started."
        break
    fi
    echo "Waiting for Grafana container to be up and running..."
    sleep 1
done

SERVICE_ACCOUNT_ID=$(curl -X POST -s -H "Content-Type: application/json" -d '{"name":"test", "role": "Admin"}' http://localhost:3000/api/serviceaccounts | jq -r .id)

if [[ -z "$SERVICE_ACCOUNT_ID" || "$SERVICE_ACCOUNT_ID" == null ]]; then
    echo "Failed to generate Grafana token"
    echo "You can either execute 'make stop && make start' to try again or execute './scripts/generate_grafana_token.sh' script manually"
    exit 1
fi

echo "Service account id: $SERVICE_ACCOUNT_ID"
GRAFANA_TOKEN=$(curl -X POST -s -H "Content-Type: application/json" -H "Authorization: Basic $(echo -n 'admin:admin' | base64)" -d "{\"name\": \"test-token-$SERVICE_ACCOUNT_ID\", \"secondsToLive\": 86400}" http://localhost:3000/api/serviceaccounts/$SERVICE_ACCOUNT_ID/tokens | jq .key)
echo "Grafana token: $GRAFANA_TOKEN"
