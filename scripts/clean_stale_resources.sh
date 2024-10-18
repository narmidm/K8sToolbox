#!/bin/bash
# clean_stale_resources.sh - Clean up stale or unused Kubernetes resources

NAMESPACE=${1:-default}

# Function to clean up completed jobs
echo "Cleaning up completed jobs in namespace: $NAMESPACE"
for JOB in $(kubectl get jobs -n $NAMESPACE --field-selector=status.succeeded=1 -o jsonpath='{.items[*].metadata.name}')
do
    echo "Deleting completed job: $JOB"
    kubectl delete job $JOB -n $NAMESPACE
done

# Function to clean up old replicasets
echo "Cleaning up old replicasets in namespace: $NAMESPACE"
for RS in $(kubectl get rs -n $NAMESPACE --no-headers | awk '$2 == 0 {print $1}')
do
    echo "Deleting old replicaset: $RS"
    kubectl delete rs $RS -n $NAMESPACE
done

# Function to clean up orphaned persistent volumes
echo "Cleaning up orphaned persistent volumes"
for PV in $(kubectl get pv --field-selector=status.phase=Released -o jsonpath='{.items[*].metadata.name}')
do
    echo "Deleting orphaned persistent volume: $PV"
    kubectl delete pv $PV
done

echo "Stale resource cleanup completed."
