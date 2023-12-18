#!/bin/bash

# Set the path to your Kubernetes Deployment YAML file
DEPLOYMENT_FILE="ex-aao.yaml"

# Kubernetes API server URL
K8S_API_SERVER="https://kubernetes.default"

# Token for authentication (try using service account token first)
TOKEN_FILE="/var/run/secrets/kubernetes.io/serviceaccount/token"
TOKEN=""

if [ -f "$TOKEN_FILE" ]; then
  TOKEN=$(cat "$TOKEN_FILE")
  echo "Using service account token for authentication."
else
  echo "Service account token not found. Make sure the script is running inside a Kubernetes Pod with the necessary permissions."
  exit 1
fi

# Apply the deployment to the cluster using cURL
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" --data-binary "@$DEPLOYMENT_FILE" "$K8S_API_SERVER/apis/apps/v1/namespaces/default/deployments"

# Check if the deployment was successful
if [ $? -eq 0 ]; then
  echo "Deployment applied successfully."
else
  echo "Error: Deployment failed."
  exit 1
fi

# Add additional checks or tests as needed
# For simplicity, we're using a sleep here, but you might want to implement a more robust wait mechanism.
sleep 30

# Check if the pods are ready using cURL or other suitable method
# ...

# Exit with success
exit 0
