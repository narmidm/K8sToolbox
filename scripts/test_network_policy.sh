#!/bin/bash
# test_network_policy.sh - Test Kubernetes network policies between pods

NAMESPACE=${1:-default}
SOURCE_POD=$2
TARGET_POD=$3

if [ -z "$SOURCE_POD" ] || [ -z "$TARGET_POD" ]; then
    echo "Usage: $0 <namespace> <source_pod> <target_pod>"
    exit 1
fi

echo "Testing network connectivity from pod $SOURCE_POD to pod $TARGET_POD in namespace: $NAMESPACE"

# Execute ping command from source pod to target pod
kubectl exec $SOURCE_POD -n $NAMESPACE -- ping -c 4 $TARGET_POD

if [ $? -eq 0 ]; then
    echo "Network connectivity test successful: $SOURCE_POD can reach $TARGET_POD."
else
    echo "Network connectivity test failed: $SOURCE_POD cannot reach $TARGET_POD."
    echo "Check network policies or any other network restrictions."
fi
