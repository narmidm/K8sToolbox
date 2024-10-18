#!/bin/bash
# auto_scaling.sh - Automatically scale Kubernetes deployments based on resource usage

DEPLOYMENT=$1
NAMESPACE=${2:-default}
CPU_THRESHOLD=75  # Scale if CPU usage exceeds 75%

if [ -z "$DEPLOYMENT" ]; then
    echo "Usage: $0 <deployment_name> [namespace]"
    exit 1
fi

echo "Checking CPU usage for deployment: $DEPLOYMENT in namespace: $NAMESPACE"

CPU_USAGE=$(kubectl top pods -n $NAMESPACE --selector=app=$DEPLOYMENT --no-headers | awk '{sum += $3} END {print int(sum / NR)}')

if [ "$CPU_USAGE" -gt "$CPU_THRESHOLD" ]; then
    echo "CPU usage $CPU_USAGE% exceeds threshold $CPU_THRESHOLD%. Scaling up deployment: $DEPLOYMENT"
    kubectl scale deployment $DEPLOYMENT -n $NAMESPACE --replicas=$(kubectl get deployment $DEPLOYMENT -n $NAMESPACE -o jsonpath='{.spec.replicas}' | awk '{print $1 + 1}')
else
    echo "CPU usage $CPU_USAGE% is within limits. No scaling action required."
fi
