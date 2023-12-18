#!/bin/bash

# Set the path to your Kubernetes Deployment YAML file
DEPLOYMENT_FILE="ex-aao.yaml"

# Apply the deployment to the cluster
kubectl apply -f "${DEPLOYMENT_FILE}"

# Check if the deployment was successful
if [ $? -eq 0 ]; then
  echo "Deployment applied successfully."
else
  echo "Error: Deployment failed."
  exit 1
fi

# Add additional checks or tests as needed, for example, wait for pods to be ready
# For simplicity, we're using a sleep here, but you might want to implement a more robust wait mechanism.
sleep 30

# Check if the pods are ready
if kubectl wait --for=condition=Ready pod --timeout=60s -l app=ex-aao > /dev/null; then
  echo "Pods are ready. Deployment test passed."
else
  echo "Error: Pods are not ready. Deployment test failed."
  exit 1
fi

# Clean up (optional): Uncomment the line below if you want to delete the deployment after testing
# kubectl delete -f "${DEPLOYMENT_FILE}"

# Exit with success
exit 0
