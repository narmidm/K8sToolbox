#!/bin/bash
# delete_stuck_namespace.sh - Delete stuck namespaces in a Kubernetes cluster

if [ -z "$1" ]; then
    echo "Usage: $0 <namespace>"
    exit 1
fi

NAMESPACE=$1
echo "Attempting to delete namespace: $NAMESPACE"

# Remove finalizers to force delete the namespace
kubectl get namespace $NAMESPACE -o json | jq '.spec = {"finalizers":[]}' | kubectl replace --raw /api/v1/namespaces/$NAMESPACE/finalize -f -

# Delete the namespace
kubectl delete namespace $NAMESPACE --grace-period=0 --force

if [ $? -eq 0 ]; then
    echo "Namespace $NAMESPACE deleted successfully."
else
    echo "Failed to delete namespace $NAMESPACE."
fi
