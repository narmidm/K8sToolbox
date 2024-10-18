#!/bin/bash
# restart_failed_pods.sh - Restart all failed pods in a Kubernetes namespace

# Set default namespace to 'default' if not provided
NAMESPACE=${1:-default}

echo "Restarting failed pods in namespace: $NAMESPACE"

for POD in $(kubectl get pods -n $NAMESPACE --field-selector=status.phase=Failed -o jsonpath='{.items[*].metadata.name}')
do
    echo "Restarting pod: $POD"
    kubectl delete pod $POD -n $NAMESPACE
done
