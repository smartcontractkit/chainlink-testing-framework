#!/bin/bash

SERVICE_ACCOUNT_ID=$(curl -X POST -s -H "Content-Type: application/json" -d '{"name":"test", "role": "Admin"}' http://localhost:3000/api/serviceaccounts | jq -r .id)
echo "Service account id: $SERVICE_ACCOUNT_ID"
GRAFANA_TOKEN=$(curl -X POST -s -H "Content-Type: application/json" -H "Authorization: Basic $(echo -n 'admin:admin' | base64)" -d "{\"name\": \"test-token-$SERVICE_ACCOUNT_ID\", \"secondsToLive\": 86400}" http://localhost:3000/api/serviceaccounts/$SERVICE_ACCOUNT_ID/tokens | jq .key)
echo "Grafana token: $GRAFANA_TOKEN"