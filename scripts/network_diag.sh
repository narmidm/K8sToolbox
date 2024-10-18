#!/bin/bash
# network_diag.sh - Advanced network diagnostics for Kubernetes pods and services

NAMESPACE=${1:-default}
TARGET_POD=$2
if [ -z "$TARGET_POD" ]; then
    echo "Usage: $0 <namespace> <target_pod>"
    exit 1
fi

echo "Starting network diagnostics for pod: $TARGET_POD in namespace: $NAMESPACE"

# Check network policies in the namespace
echo "Checking network policies in namespace: $NAMESPACE"
kubectl get networkpolicies -n $NAMESPACE

# Verify endpoints for services in the namespace
echo "Checking service endpoints in namespace: $NAMESPACE"
kubectl get endpoints -n $NAMESPACE

# Running tcpdump to capture traffic for analysis
echo "Capturing network traffic on pod: $TARGET_POD"
kubectl exec $TARGET_POD -n $NAMESPACE -- tcpdump -c 20 -w /tmp/traffic.pcap
echo "Traffic captured in /tmp/traffic.pcap inside pod: $TARGET_POD"
