[![CI Status](https://github.com/narmidm/k8stoolbox/actions/workflows/ci.yml/badge.svg)](https://github.com/narmidm/K8sToolbox/actions/workflows/ci.yml)
[![CD Status](https://github.com/narmidm/k8stoolbox/actions/workflows/cd.yml/badge.svg)](https://github.com/narmidm/K8sToolbox/actions/workflows/cd.yml)
[![Docker Image Version](https://img.shields.io/docker/v/narmidm/k8stoolbox?sort=semver)](https://hub.docker.com/repository/docker/narmidm/k8stoolbox)
[![Docker Pulls](https://img.shields.io/docker/pulls/narmidm/k8stoolbox)](https://hub.docker.com/repository/docker/narmidm/k8stoolbox)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/narmidm/K8sToolbox)](https://raw.githubusercontent.com/narmidm/K8sToolbox/refs/heads/master/go.mod)
[![GitHub License](https://img.shields.io/github/license/narmidm/K8sToolbox)](https://raw.githubusercontent.com/narmidm/K8sToolbox/refs/heads/master/LICENSE)
[![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/that_imran)](https://x.com/intent/user?screen_name=that_imran)
[![Linktree](https://img.shields.io/badge/linktree-1de9b6?style=for-the-badge&logo=linktree&logoColor=white)](https://linktr.ee/narmidm)
<a href="https://www.linkedin.com/comm/mynetwork/discovery-see-all?usecase=PEOPLE_FOLLOWS&followMember=narmidm" target="blank"><img src="https://img.shields.io/badge/LinkedIn-Connect-blue" alt="narmidm" /></a>
![Contributors](https://img.shields.io/github/contributors/narmidm/k8stoolbox)
[![GitHub Issues](https://img.shields.io/github/issues/narmidm/k8stoolbox)](https://github.com/narmidm/K8sToolbox/issues)
[![GitHub Stars](https://img.shields.io/github/stars/narmidm/k8stoolbox)](https://github.com/narmidm/K8sToolbox/stargazers)
[![GitHub Forks](https://img.shields.io/github/forks/narmidm/k8stoolbox)](https://github.com/narmidm/K8sToolbox/forks)
[![Last Commit](https://img.shields.io/github/last-commit/narmidm/k8stoolbox)](https://github.com/narmidm/K8sToolbox/commits/master/)



# K8sToolbox
<img src="assets/logo.png" alt="K8sToolbox Logo" width="160"/>

**K8sToolbox** is a versatile toolkit designed for managing Kubernetes clusters, providing essential debugging tools and utilities for cluster administrators and developers.

## Overview

K8sToolbox is an all-in-one, versatile toolkit engineered to streamline the management and troubleshooting of Kubernetes clusters. It serves as a one-stop solution for all Kubernetes management and troubleshooting needs, providing everything required to maintain and optimize your cluster environments. It provides a comprehensive suite of powerful debugging, diagnostic, and operational tools designed specifically for cluster administrators, DevOps engineers, and developers. With K8sToolbox, users gain the ability to efficiently manage workloads, diagnose issues, and maintain optimal performance of Kubernetes environments.

Equipped with a rich collection of utilities that include advanced network diagnostics, automated health checks, resource monitoring, and recovery mechanisms, K8sToolbox empowers teams to proactively identify and resolve cluster-related challenges. Whether you need to aggregate logs, test pod connectivity, manage resources, or troubleshoot network policies, K8sToolbox offers a unified solution that significantly reduces complexity while enhancing productivity.

K8sToolbox also integrates essential third-party utilities like `kubectl`, `stern`, and `k9s`, providing a seamless, command-driven experience for interacting with Kubernetes clusters. This toolkit not only simplifies debugging processes but also provides the scalability and robustness necessary to maintain large-scale, multi-node environments effectively. By combining a vast array of capabilities into a single, easy-to-use image, K8sToolbox ensures that Kubernetes management is more accessible, efficient, and reliable than ever before.

## Use Cases
**K8sToolbox** is perfect for:
- **Cluster Troubleshooting**: Quickly diagnose issues in your cluster, such as resource contention, network issues, or failed pods.
- **Maintenance**: Clean up stale or unused resources like completed jobs and old replicasets to keep your cluster healthy.
- **Automation**: Automate tasks like scaling deployments, resource usage checks, and more.
- **Debugging Network Policies**: Validate network connectivity and ensure your network policies are properly configured.
- **Log Aggregation**: Collect and analyze logs from multiple namespaces and pods to understand the state of your cluster and applications.

With **K8sToolbox**, you can:
- Execute health checks, manage stuck resources, aggregate logs, and perform network diagnostics.
- Run custom scripts directly from your local machine or inside a Kubernetes pod using shell exec. **"For detailed usage instructions on the scripts, please refer to [Using K8sToolbox Scripts](#using-k8stoolbox-scripts)."**


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
kubectl apply -f https://raw.githubusercontent.com/narmidm/K8sToolbox/refs/heads/master/manifests/debug-pod.yaml
```

This creates a pod named `k8stoolbox-debug` in the `default` namespace, which can be used for one-off debugging and troubleshooting tasks.

#### DaemonSet
To deploy **K8sToolbox** on all nodes, use the DaemonSet manifest:

```sh
kubectl apply -f https://raw.githubusercontent.com/narmidm/K8sToolbox/refs/heads/master/manifests/debug-daemon.yaml
```

This creates a **DaemonSet** that runs **K8sToolbox** on all nodes, making it accessible from anywhere in the cluster.

### Utilizing K8sToolbox
There are two primary ways to use **K8sToolbox**:
1. **Local Execution**: Running scripts directly from the local system.
2. **Kubernetes Shell Execution**: Executing commands inside a running **K8sToolbox** pod using `kubectl exec`.

#### 1. Running Scripts Locally
You can run the scripts in the `/scripts` directory locally if you have **kubectl** configured and connected to your Kubernetes cluster.
**"For detailed usage instructions on the scripts, please refer to [Using K8sToolbox Scripts](#using-k8stoolbox-scripts)."**


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

## Using K8sToolbox Scripts

This section provides detailed information on the various utility scripts available in the `K8sToolbox` repository, located under the `/scripts` directory. Each script has been crafted to assist with common Kubernetes cluster management and troubleshooting tasks. Below, you will find how to use each script, including commands to run them directly or through the symlinks set up during installation.

To learn how to use these scripts, follow the instructions provided for each script below. You can execute these commands from your local machine or within a Kubernetes pod, depending on your setup.

### Running the Scripts Locally
If you have cloned the repository and have `kubectl` configured to interact with your cluster, you can run the scripts directly by using:

```sh
./scripts/<script_name>.sh [arguments]
```

Alternatively, you can use the symlinks created for each script, allowing you to run them without specifying the full path:

```sh
<script_name> [arguments]
```

Ensure that your shell environment includes `/usr/local/bin` in the `$PATH` so that the symlinks are accessible.

### Script Details
Below is a list of all available scripts, with detailed descriptions and examples of how to use them:

1. **aggregate_logs.sh**  
   Aggregates logs from all pods within a specified namespace. This script is useful when you need a combined view of application logs.
   ```sh
   aggregate_logs <namespace>
   ```
   Example:
   ```sh
   aggregate_logs default
   ```
   This command will aggregate logs from all pods in the `default` namespace and print them to the console.

2. **auto_recover.sh**  
   Automatically recovers failed pods and restarts them as needed. This script can be run to automate pod recovery.
   ```sh
   auto_recover <namespace>
   ```
   Example:
   ```sh
   auto_recover kube-system
   ```
   This command will automatically recover any failed pods in the `kube-system` namespace by restarting them.

3. **auto_scaling.sh**  
   Adjusts deployment scaling based on resource usage. This script can help automate horizontal scaling.
   ```sh
   auto_scaling <namespace> <deployment_name> <desired_replicas>
   ```
   Example:
   ```sh
   auto_scaling default my-app 5
   ```
   This command will scale the deployment named `my-app` in the `default` namespace to `5` replicas.

4. **backup_restore.sh**  
   Backs up and restores Kubernetes resources in a namespace. Useful for disaster recovery scenarios.
   ```sh
   backup_restore backup <namespace>
   backup_restore restore <namespace>
   ```
   Example:
   ```sh
   backup_restore backup default
   ```
   This command will back up all resources in the `default` namespace. Use `restore` instead of `backup` to restore the resources.

5. **clean_stale_resources.sh**  
   Cleans up old Kubernetes resources such as completed jobs and replicasets.
   ```sh
   clean_stale_resources <namespace>
   ```
   Example:
   ```sh
   clean_stale_resources default
   ```
   This command will clean up stale resources (e.g., completed jobs, old replicasets) in the `default` namespace.

6. **connectivity_test.sh**  
   Tests network connectivity between two pods. Useful for validating network policies.
   ```sh
   connectivity_test <namespace> <source_pod> <target_pod>
   ```
   Example:
   ```sh
   connectivity_test default pod-a pod-b
   ```
   This command will test network connectivity from `pod-a` to `pod-b` in the `default` namespace.

7. **delete_stuck_crds.sh**  
   Deletes Custom Resource Definitions (CRDs) that are stuck due to finalizers.
   ```sh
   delete_stuck_crds <crd_name>
   ```
   Example:
   ```sh
   delete_stuck_crds my-crd
   ```
   This command will forcefully delete the CRD named `my-crd` that is stuck due to finalizers.

8. **delete_stuck_namespace.sh**  
   Deletes namespaces that are stuck in terminating status.
   ```sh
   delete_stuck_namespace <namespace>
   ```
   Example:
   ```sh
   delete_stuck_namespace default
   ```
   This command will delete the `default` namespace if it is stuck in a terminating state.

9. **healthcheck.sh**  
   Performs health checks on all pods within a namespace.
   ```sh
   healthcheck <namespace>
   ```
   Example:
   ```sh
   healthcheck kube-system
   ```
   This command will perform health checks on all pods in the `kube-system` namespace and report any issues found.

10. **network_diag.sh**  
    Provides advanced network diagnostics, including traffic capture between pods.
    ```sh
    network_diag <namespace> <source_pod> <target_pod>
    ```
    Example:
    ```sh
    network_diag default pod-a pod-b
    ```
    This command will perform network diagnostics between `pod-a` and `pod-b` in the `default` namespace, including capturing traffic if needed.

11. **resource_usage.sh**  
    Monitors CPU and memory usage for nodes and pods within a namespace.
    ```sh
    resource_usage <namespace>
    ```
    Example:
    ```sh
    resource_usage default
    ```
    This command will monitor and display the CPU and memory usage for all nodes and pods in the `default` namespace.

12. **restart_failed_pods.sh**  
    Restarts all failed pods in a given namespace.
    ```sh
    restart_failed_pods <namespace>
    ```
    Example:
    ```sh
    restart_failed_pods kube-system
    ```
    This command will restart all failed pods in the `kube-system` namespace.

13. **snapshot_audit.sh**  
    Takes a snapshot of the cluster state for auditing purposes.
    ```sh
    snapshot_audit <namespace>
    ```
    Example:
    ```sh
    snapshot_audit default
    ```
    This command will take a snapshot of the cluster state for the `default` namespace, which can be used for auditing purposes.

14. **test_network_policy.sh**  
    Tests network policies by attempting connections between source and target pods.
    ```sh
    test_network_policy <namespace> <source_pod> <target_pod>
    ```
    Example:
    ```sh
    test_network_policy default pod-a pod-b
    ```
    This command will test the network policies by attempting to connect from `pod-a` to `pod-b` in the `default` namespace.

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
