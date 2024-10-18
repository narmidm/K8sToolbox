#!/bin/bash
# connectivity_test.sh - Test network connectivity between pods or services in a Kubernetes cluster

# Check if required arguments are provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <source_pod_name> <target_service_or_pod_ip>"
    exit 1
fi

SOURCE_POD=$1
TARGET=$2

echo "Testing connectivity from pod $SOURCE_POD to target $TARGET"
kubectl exec $SOURCE_POD -- ping -c 4 $TARGET
