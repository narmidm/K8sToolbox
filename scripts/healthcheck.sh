#!/bin/bash
# healthcheck.sh - A script to perform comprehensive health checks in a Kubernetes cluster
# 
# This script checks the health status of pods, services, and nodes in a specified namespace.
# It collects detailed metrics and identifies potential issues.

set -eo pipefail

# Script version
VERSION="1.0.0"

# Terminal colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Variables
NAMESPACE=""
VERBOSE=false
TIMEOUT=30
OUTPUT_FORMAT="text" # text, json, yaml

# Function to display usage information
function display_usage() {
    echo -e "${BLUE}K8sToolbox - Health Check Utility v${VERSION}${NC}"
    echo
    echo "This script performs health checks on Kubernetes resources."
    echo
    echo -e "${YELLOW}Usage:${NC}"
    echo "  $0 [options] [namespace]"
    echo
    echo -e "${YELLOW}Options:${NC}"
    echo "  -h, --help                Display this help message"
    echo "  -v, --verbose             Enable verbose output"
    echo "  -t, --timeout SECONDS     Set timeout for operations (default: 30s)"
    echo "  -o, --output FORMAT       Output format: text, json, yaml (default: text)"
    echo
    echo -e "${YELLOW}Examples:${NC}"
    echo "  $0 default                 Check health in default namespace"
    echo "  $0 -v kube-system          Verbose health check in kube-system namespace"
    echo "  $0 --output json default   Output results in JSON format"
}

# Function to log messages
function log() {
    local level=$1
    local message=$2
    local timestamp=$(date "+%Y-%m-%d %H:%M:%S")
    
    case $level in
        "INFO")
            echo -e "${GREEN}[INFO]${NC} ${timestamp} - ${message}"
            ;;
        "WARN")
            echo -e "${YELLOW}[WARN]${NC} ${timestamp} - ${message}"
            ;;
        "ERROR")
            echo -e "${RED}[ERROR]${NC} ${timestamp} - ${message}"
            ;;
        *)
            echo -e "${BLUE}[DEBUG]${NC} ${timestamp} - ${message}"
            ;;
    esac
}

# Function to check if a command exists
function check_command() {
    if ! command -v "$1" &> /dev/null; then
        log "ERROR" "$1 could not be found. Please install $1 to use this script."
        exit 1
    fi
}

# Function to check pod health
function check_pod_health() {
    local ns=$1
    local pod_name=$2
    local issues=()
    
    # Get pod status
    local status=$(kubectl get pod "$pod_name" -n "$ns" -o jsonpath='{.status.phase}' 2>/dev/null)
    if [ -z "$status" ]; then
        log "ERROR" "Failed to get status for pod $pod_name"
        return 1
    fi
    
    # Check for container restarts
    local restarts=$(kubectl get pod "$pod_name" -n "$ns" -o jsonpath='{.status.containerStatuses[*].restartCount}' 2>/dev/null)
    if [ ! -z "$restarts" ]; then
        for restart in $restarts; do
            if [ "$restart" -gt 5 ]; then
                issues+=("High restart count: $restart")
            fi
        done
    fi
    
    # Check ready status
    local ready=$(kubectl get pod "$pod_name" -n "$ns" -o jsonpath='{.status.containerStatuses[*].ready}' 2>/dev/null)
    if [[ $ready == *"false"* ]]; then
        issues+=("Container not ready")
    fi
    
    # Determine overall health
    if [ "$status" != "Running" ] && [ "$status" != "Succeeded" ]; then
        issues+=("Status: $status")
    fi
    
    # Report results
    if [ ${#issues[@]} -eq 0 ]; then
        log "INFO" "Pod $pod_name is healthy"
        return 0
    else
        log "WARN" "Pod $pod_name has issues:"
        for issue in "${issues[@]}"; do
            log "WARN" "  - $issue"
        done
        return 1
    fi
}

# Function to verify endpoint connectivity
function check_service_endpoints() {
    local ns=$1
    local svc=$2
    
    # Get endpoints
    local endpoints=$(kubectl get endpoints "$svc" -n "$ns" -o jsonpath='{.subsets[*].addresses[*].ip}' 2>/dev/null)
    
    if [ -z "$endpoints" ]; then
        log "WARN" "No endpoints found for service $svc"
        return 1
    else
        if [ "$VERBOSE" = true ]; then
            log "INFO" "Service $svc has endpoints: $endpoints"
        else
            log "INFO" "Service $svc has $(echo $endpoints | wc -w) endpoint(s)"
        fi
        return 0
    fi
}

# Function to check node conditions
function check_node_conditions() {
    local nodes=$(kubectl get nodes -o jsonpath='{.items[*].metadata.name}' 2>/dev/null)
    local healthy_nodes=0
    local total_nodes=0
    
    log "INFO" "Checking node conditions..."
    
    for node in $nodes; do
        total_nodes=$((total_nodes + 1))
        local conditions=$(kubectl get node "$node" -o jsonpath='{.status.conditions[?(@.status=="True")].type}' 2>/dev/null)
        local unhealthy=false
        
        if [[ ! $conditions == *"Ready"* ]]; then
            log "WARN" "Node $node is not Ready"
            unhealthy=true
        fi
        
        if [[ $conditions == *"MemoryPressure"* ]]; then
            log "WARN" "Node $node is experiencing Memory Pressure"
            unhealthy=true
        fi
        
        if [[ $conditions == *"DiskPressure"* ]]; then
            log "WARN" "Node $node is experiencing Disk Pressure"
            unhealthy=true
        fi
        
        if [[ $conditions == *"PIDPressure"* ]]; then
            log "WARN" "Node $node is experiencing PID Pressure"
            unhealthy=true
        fi
        
        if [ "$unhealthy" = false ]; then
            healthy_nodes=$((healthy_nodes + 1))
            if [ "$VERBOSE" = true ]; then
                log "INFO" "Node $node is healthy"
            fi
        fi
    done
    
    log "INFO" "$healthy_nodes/$total_nodes nodes are healthy"
}

# Main health check function
function run_health_check() {
    local ns=$1
    local timestart=$(date +%s)
    
    log "INFO" "Starting health check for namespace: $ns"
    
    # Validate namespace exists
    if ! kubectl get namespace "$ns" &>/dev/null; then
        log "ERROR" "Namespace $ns does not exist"
        exit 1
    fi
    
    # Check pod health
    log "INFO" "Checking pod status in namespace: $ns"
    local pods=$(kubectl get pods -n "$ns" -o jsonpath='{.items[*].metadata.name}' 2>/dev/null)
    local healthy_pods=0
    local total_pods=0
    
    if [ -z "$pods" ]; then
        log "INFO" "No pods found in namespace $ns"
    else
        for pod in $pods; do
            total_pods=$((total_pods + 1))
            if check_pod_health "$ns" "$pod"; then
                healthy_pods=$((healthy_pods + 1))
            fi
        done
        
        log "INFO" "$healthy_pods/$total_pods pods are healthy in namespace $ns"
    fi
    
    # Check services
    log "INFO" "Checking services in namespace: $ns"
    local services=$(kubectl get services -n "$ns" -o jsonpath='{.items[*].metadata.name}' 2>/dev/null)
    local healthy_services=0
    local total_services=0
    
    if [ -z "$services" ]; then
        log "INFO" "No services found in namespace $ns"
    else
        for svc in $services; do
            total_services=$((total_services + 1))
            if check_service_endpoints "$ns" "$svc"; then
                healthy_services=$((healthy_services + 1))
            fi
        done
        
        log "INFO" "$healthy_services/$total_services services have endpoints in namespace $ns"
    fi
    
    # Check nodes (if verbose)
    if [ "$VERBOSE" = true ]; then
        check_node_conditions
    fi
    
    # Output summary
    local timeend=$(date +%s)
    local duration=$((timeend - timestart))
    
    log "INFO" "Health check completed in $duration seconds"
    if [ $healthy_pods -eq $total_pods ] && [ $healthy_services -eq $total_services ]; then
        log "INFO" "All resources in namespace $ns are healthy"
        return 0
    else
        log "WARN" "Some resources in namespace $ns have issues"
        return 1
    fi
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            display_usage
            exit 0
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -t|--timeout)
            TIMEOUT="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_FORMAT="$2"
            if [[ ! "$OUTPUT_FORMAT" =~ ^(text|json|yaml)$ ]]; then
                log "ERROR" "Invalid output format: $OUTPUT_FORMAT"
                display_usage
                exit 1
            fi
            shift 2
            ;;
        *)
            if [ -z "$NAMESPACE" ]; then
                NAMESPACE="$1"
            else
                log "ERROR" "Unknown argument: $1"
                display_usage
                exit 1
            fi
            shift
            ;;
    esac
done

# Set default namespace if not provided
if [ -z "$NAMESPACE" ]; then
    NAMESPACE="default"
fi

# Check if required tools are installed
check_command "kubectl"

# Set a timeout for long-running kubectl commands
export KUBECTL_TIMEOUT="${TIMEOUT}s"

# Run health check
run_health_check "$NAMESPACE"
exit_code=$?

# Exit with appropriate code
exit $exit_code
