#!/bin/bash
# snapshot_audit.sh - Take a snapshot of the current Kubernetes cluster state for auditing

OUTPUT_DIR="./cluster_snapshot_$(date +%Y%m%d_%H%M%S)"
mkdir -p $OUTPUT_DIR

echo "Taking a snapshot of the Kubernetes cluster state..."

# Collect all namespaces
echo "Collecting namespace data..."
kubectl get namespaces -o yaml > $OUTPUT_DIR/namespaces.yaml

# Collect all pods, services, deployments, configmaps, etc.
for RESOURCE in pods services deployments configmaps nodes
do
    echo "Collecting data for resource: $RESOURCE"
    kubectl get $RESOURCE --all-namespaces -o yaml > $OUTPUT_DIR/${RESOURCE}.yaml
done

# Collect node resource utilization
echo "Collecting node resource usage..."
kubectl top nodes > $OUTPUT_DIR/node_usage.txt

# Collect cluster events
echo "Collecting cluster events..."
kubectl get events --all-namespaces > $OUTPUT_DIR/cluster_events.txt

echo "Cluster state snapshot saved to: $OUTPUT_DIR"
