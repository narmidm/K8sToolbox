#!/bin/bash
# delete_stuck_crds.sh - Delete stuck CRDs in a Kubernetes cluster

if [ -z "$1" ]; then
    echo "Usage: $0 <crd_name>"
    exit 1
fi

CRD=$1
echo "Attempting to delete CRD: $CRD"

# Patch CRD to remove finalizers if it is stuck
echo "Patching CRD to remove finalizers..."
kubectl patch crd $CRD -p '{"metadata":{"finalizers":null}}' --type=merge

# Delete all resources associated with the CRD
echo "Deleting all instances of the CRD: $CRD"
for RESOURCE in $(kubectl get $CRD --all-namespaces -o jsonpath='{.items[*].metadata.name}')
do
    kubectl delete $CRD $RESOURCE --all-namespaces --force --grace-period=0
    if [ $? -ne 0 ]; then
        echo "Failed to delete resource: $RESOURCE"
    else
        echo "Resource $RESOURCE deleted successfully."
    fi
done

# Delete the CRD itself
echo "Deleting the CRD: $CRD"
kubectl delete crd $CRD --grace-period=0 --force
if [ $? -eq 0 ]; then
    echo "CRD $CRD deleted successfully."
else
    echo "Failed to delete CRD $CRD."
fi
