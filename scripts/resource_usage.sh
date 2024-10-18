#!/bin/bash
# resource_usage.sh - A script to monitor CPU and memory usage for Kubernetes nodes and pods

echo "Checking node resource usage..."
kubectl top nodes

echo "Checking pod resource usage..."
kubectl top pods --all-namespaces
