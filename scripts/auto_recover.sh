#!/bin/bash
# auto_recover.sh - Automatically recover non-running pods and send an alert

NAMESPACE=${1:-default}
ALERT_WEBHOOK_URL="https://your-webhook-url.com/alert"

echo "Checking and recovering non-running pods in namespace: $NAMESPACE"

for POD in $(kubectl get pods -n $NAMESPACE --field-selector=status.phase!=Running -o jsonpath='{.items[*].metadata.name}')
do
    echo "Pod $POD is not running. Attempting recovery..."
    kubectl delete pod $POD -n $NAMESPACE
    sleep 5  # Wait a bit before checking if it comes back up

    # Check pod status again
    NEW_STATUS=$(kubectl get pod $POD -n $NAMESPACE -o jsonpath='{.status.phase}' 2>/dev/null)
    if [ "$NEW_STATUS" != "Running" ]; then
        echo "Pod $POD recovery failed. Sending alert..."
        curl -X POST -H 'Content-type: application/json' --data '{"text":"Pod recovery failed for pod: '"$POD"' in namespace: '"$NAMESPACE"'."}' $ALERT_WEBHOOK_URL
    else
        echo "Pod $POD successfully recovered."
    fi
    echo "------------------------------------"
done
