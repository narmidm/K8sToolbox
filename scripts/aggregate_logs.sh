#!/bin/bash
# aggregate_logs.sh - Collect logs from multiple namespaces and pods for analysis

NAMESPACES=${@:-default}
LOG_DIR="./logs_$(date +%Y%m%d_%H%M%S)"

mkdir -p $LOG_DIR

echo "Aggregating logs for namespaces: $NAMESPACES"

for NAMESPACE in $NAMESPACES
 do
    echo "Collecting logs from namespace: $NAMESPACE"
    for POD in $(kubectl get pods -n $NAMESPACE -o jsonpath='{.items[*].metadata.name}')
    do
        echo "Collecting logs for pod: $POD"
        kubectl logs $POD -n $NAMESPACE > $LOG_DIR/${NAMESPACE}_${POD}.log
    done
    echo "------------------------------------"
done

echo "Logs aggregated and saved in directory: $LOG_DIR"
