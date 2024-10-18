#!/bin/bash
# healthcheck.sh - A script to perform basic health checks in a Kubernetes cluster

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null
then
    echo "kubectl could not be found. Please install kubectl to use this script."
    exit 1
fi

# Set default namespace to 'default' if not provided
NAMESPACE=${1:-default}

# List all pods in the namespace and show their status
echo "Checking pod status in namespace: $NAMESPACE"
kubectl get pods -n $NAMESPACE

# Perform detailed health checks on each pod
for POD in $(kubectl get pods -n $NAMESPACE -o jsonpath='{.items[*].metadata.name}')
do
    STATUS=$(kubectl get pod $POD -n $NAMESPACE -o jsonpath='{.status.phase}')
    echo "Pod: $POD Status: $STATUS"

    if [ "$STATUS" != "Running" ]; then
        echo "Fetching logs for pod: $POD"
        kubectl logs $POD -n $NAMESPACE --tail=10
    fi
	echo "------------------------------------"
done

# Check node status
echo "Checking node statuses..."
kubectl get nodes -o wide

# Verify services in the namespace
echo "Checking services in namespace: $NAMESPACE"
kubectl get svc -n $NAMESPACE
