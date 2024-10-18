[![CI Status](https://github.com/narmidm/k8stoolbox/actions/workflows/ci.yml/badge.svg)](https://github.com/narmidm/K8sToolbox/actions/workflows/ci.yml)
[![CD Status](https://github.com/narmidm/k8stoolbox/actions/workflows/cd.yml/badge.svg)](https://github.com/narmidm/K8sToolbox/actions/workflows/cd.yml)
[![Docker Image Version](https://img.shields.io/docker/v/narmidm/k8stoolbox?sort=semver)](https://hub.docker.com/repository/docker/narmidm/k8stoolbox)
[![Docker Pulls](https://img.shields.io/docker/pulls/narmidm/k8stoolbox)](https://hub.docker.com/repository/docker/narmidm/k8stoolbox)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/narmidm/K8sToolbox)](https://raw.githubusercontent.com/narmidm/K8sToolbox/refs/heads/master/go.mod)
[![GitHub License](https://img.shields.io/github/license/narmidm/K8sToolbox)](https://raw.githubusercontent.com/narmidm/K8sToolbox/refs/heads/master/LICENSE)
[![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/that_imran)](https://x.com/that_imran)
<a href="https://www.linkedin.com/comm/mynetwork/discovery-see-all?usecase=PEOPLE_FOLLOWS&followMember=narmidm" target="blank"><img src="https://img.shields.io/badge/LinkedIn-Connect-blue" alt="narmidm" /></a>
![Contributors](https://img.shields.io/github/contributors/narmidm/k8stoolbox)
[![GitHub Issues](https://img.shields.io/github/issues/narmidm/k8stoolbox)](https://github.com/narmidm/K8sToolbox/issues)
[![GitHub Stars](https://img.shields.io/github/stars/narmidm/k8stoolbox)](https://github.com/narmidm/K8sToolbox/stargazers)
[![GitHub Forks](https://img.shields.io/github/forks/narmidm/k8stoolbox)](https://github.com/narmidm/K8sToolbox/forks)
[![Last Commit](https://img.shields.io/github/last-commit/narmidm/k8stoolbox)](https://github.com/narmidm/K8sToolbox/commits/master/)



# K8sToolbox

### Use Cases
**K8sToolbox** is perfect for:
- **Cluster Troubleshooting**: Quickly diagnose issues in your cluster, such as resource contention, network issues, or failed pods.
- **Maintenance**: Clean up stale or unused resources like completed jobs and old replicasets to keep your cluster healthy.
- **Automation**: Automate tasks like scaling deployments, resource usage checks, and more.
- **Debugging Network Policies**: Validate network connectivity and ensure your network policies are properly configured.
- **Log Aggregation**: Collect and analyze logs from multiple namespaces and pods to understand the state of your cluster and applications.

With **K8sToolbox**, you can:
- Execute health checks, manage stuck resources, aggregate logs, and perform network diagnostics.
- Run custom scripts directly from your local machine or inside a Kubernetes pod using shell exec.

## Folder Structure
```
K8sToolbox/
│
├── docker/
│   └── Dockerfile              # Docker image definition for building K8sToolbox
│
├── manifests/
│   ├── debug-daemon.yaml       # DaemonSet manifest for deploying K8sToolbox on all nodes
│   └── debug-pod.yaml          # Pod manifest for running a standalone K8sToolbox instance
│
├── scripts/                    # Collection of helpful Kubernetes management scripts
│   ├── aggregate_logs.sh
│   ├── auto_recover.sh
│   ├── auto_scaling.sh
│   ├── backup_restore.sh
│   ├── clean_stale_resources.sh
│   ├── connectivity_test.sh
│   ├── delete_stuck_crds.sh
│   ├── delete_stuck_namespace.sh
│   ├── healthcheck.sh
│   ├── network_diag.sh
│   ├── resource_usage.sh
│   ├── restart_failed_pods.sh
│   ├── snapshot_audit.sh
│   └── test_network_policy.sh
│
├── .gitignore
├── CONTRIBUTING.md             # Guidelines for contributing to K8sToolbox
├── LICENSE                     # License details (Apache License 2.0)
├── go.mod                      # Go module definition
├── main.go                     # Main Golang utility file for K8sToolbox
└── README.md                   # Documentation (you're reading this!)
```

## Getting Started
### Prerequisites
- **Docker** installed to build and run the K8sToolbox Docker image.
- **Kubernetes cluster** with `kubectl` configured to interact with the cluster.
- **Permissions**: Ensure you have sufficient permissions to run commands like `kubectl exec` and `kubectl apply`.

### Building the Docker Image
To build the Docker image for **K8sToolbox**, run the following command in the root directory of the project:

```sh
docker build -t k8stoolbox:latest -f docker/Dockerfile .
```

This will create a Docker image named `k8stoolbox` that you can use locally or push to a container registry.

### Deploying K8sToolbox in Kubernetes
You can deploy **K8sToolbox** as either a standalone **Pod** or as a **DaemonSet** to cover all nodes.

#### Standalone Pod
To deploy a standalone **K8sToolbox** pod, use the following command:

```sh
kubectl apply -f manifests/debug-pod.yaml
```

This creates a pod named `k8stoolbox-debug` in the `default` namespace, which can be used for one-off debugging and troubleshooting tasks.

#### DaemonSet
To deploy **K8sToolbox** on all nodes, use the DaemonSet manifest:

```sh
kubectl apply -f manifests/debug-daemon.yaml
```

This creates a **DaemonSet** that runs **K8sToolbox** on all nodes, making it accessible from anywhere in the cluster.

### Utilizing K8sToolbox
There are two primary ways to use **K8sToolbox**:
1. **Local Execution**: Running scripts directly from the local system.
2. **Kubernetes Shell Execution**: Executing commands inside a running **K8sToolbox** pod using `kubectl exec`.

#### 1. Running Scripts Locally
You can run the scripts in the `/scripts` directory locally if you have **kubectl** configured and connected to your Kubernetes cluster.

Examples:

- **Backup and Restore Resources**:
  ```sh
  ./scripts/backup_restore.sh backup default
  ./scripts/backup_restore.sh restore default
  ```
- **Clean Stale Resources**:
  ```sh
  ./scripts/clean_stale_resources.sh default
  ```
- **Test Network Policies**:
  ```sh
  ./scripts/test_network_policy.sh default <source_pod> <target_pod>
  ```
- **Aggregate Logs**:
  ```sh
  ./scripts/aggregate_logs.sh default kube-system
  ```

#### 2. Running Scripts in Kubernetes Pods
You can also execute the scripts inside a running **K8sToolbox** pod by using **kubectl exec**. This is useful when you need to troubleshoot issues within the cluster itself.

First, find the name of the **K8sToolbox** pod:

```sh
kubectl get pods -n default -l app=k8stoolbox
```

Then use `kubectl exec` to run commands:

- **Execute a Health Check**:
  ```sh
  kubectl exec -it <k8stoolbox-pod-name> -- /usr/local/bin/healthcheck default
  ```
- **Run Resource Cleanup**:
  ```sh
  kubectl exec -it <k8stoolbox-pod-name> -- /usr/local/bin/clean_stale_resources default
  ```
- **Ping Between Pods to Test Network Policy**:
  ```sh
  kubectl exec -it <k8stoolbox-pod-name> -- /usr/local/bin/test_network_policy default <source_pod> <target_pod>
  ```

### Available Scripts
The `/scripts` directory contains several useful scripts for Kubernetes management:

- **aggregate_logs.sh**: Aggregates logs from all pods in specified namespaces.
- **auto_recover.sh**: Automatically recovers failed pods and sends alerts.
- **auto_scaling.sh**: Automatically scales deployments based on resource usage.
- **backup_restore.sh**: Backs up and restores Kubernetes resources in a specified namespace.
- **clean_stale_resources.sh**: Cleans up completed jobs, old replicasets, and orphaned persistent volumes.
- **connectivity_test.sh**: Tests network connectivity between pods or services.
- **delete_stuck_crds.sh**: Deletes CRDs that are stuck by removing finalizers.
- **delete_stuck_namespace.sh**: Deletes namespaces that are stuck due to finalizers.
- **healthcheck.sh**: Performs health checks on pods and nodes in a namespace.
- **network_diag.sh**: Provides advanced network diagnostics, including capturing traffic.
- **resource_usage.sh**: Monitors CPU and memory usage for nodes and pods.
- **restart_failed_pods.sh**: Restarts all failed pods in a given namespace.
- **snapshot_audit.sh**: Takes a snapshot of the cluster state for auditing purposes.
- **test_network_policy.sh**: Tests network connectivity between pods to validate network policies.

### Symlinked Commands
For convenience, all scripts are symlinked to `/usr/local/bin` in the Docker image, allowing you to call them without specifying the full path. For example:

```sh
auto_recover
backup_restore backup default
clean_stale_resources default
```

### Contributing
We welcome contributions! Please read the [CONTRIBUTING.md](CONTRIBUTING.md) file for guidelines on how to contribute to **K8sToolbox**.

### License
This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

### Authors and Acknowledgements
K8sToolbox was inspired by various Kubernetes utility tools, including the **Swiss Army Knife** for DevOps. Special thanks to all contributors who helped improve this toolbox.

### Future Work
- **Add more advanced diagnostics** tools.
- **Integration with Prometheus** for enhanced monitoring capabilities.

## Inspiration
K8sToolbox was inspired by the [swiss-army-knife](https://github.com/leodotcloud/swiss-army-knife) repository, which serves as a useful multi-purpose tool for DevOps. Our goal is to build upon that foundation and create a specialized, Kubernetes-focused toolkit that helps users effectively troubleshoot and manage their clusters.

Let us know if you have feature requests or suggestions to make **K8sToolbox** even better!
