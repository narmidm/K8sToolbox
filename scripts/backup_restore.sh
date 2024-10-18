#!/bin/bash
# backup_restore.sh - Backup and restore Kubernetes resources

ACTION=$1
NAMESPACE=${2:-default}
BACKUP_DIR="./backup_$NAMESPACE"

if [ -z "$ACTION" ]; then
    echo "Usage: $0 <backup|restore> [namespace]"
    exit 1
fi

# Function to backup resources
backup() {
    echo "Backing up resources in namespace: $NAMESPACE"
    mkdir -p $BACKUP_DIR

    for RESOURCE in pods services deployments configmaps secrets
    do
        echo "Backing up $RESOURCE..."
        kubectl get $RESOURCE -n $NAMESPACE -o yaml > $BACKUP_DIR/${RESOURCE}.yaml
    done

    echo "Backup completed and saved in directory: $BACKUP_DIR"
}

# Function to restore resources
restore() {
    if [ ! -d "$BACKUP_DIR" ]; then
        echo "Backup directory $BACKUP_DIR does not exist. Please run a backup first."
        exit 1
    fi

    echo "Restoring resources in namespace: $NAMESPACE"
    for FILE in $BACKUP_DIR/*.yaml
    do
        echo "Restoring resource from $FILE..."
        kubectl apply -f $FILE
    done

    echo "Restore completed."
}

case $ACTION in
    backup)
        backup
        ;;
    restore)
        restore
        ;;
    *)
        echo "Invalid action: $ACTION. Use 'backup' or 'restore'."
        exit 1
        ;;
esac
