# K8sToolbox

**Your Swiss Army Knife for Kubernetes Troubleshooting**

K8sToolbox is a versatile collection of tools for debugging, troubleshooting, and managing Kubernetes clusters, all bundled into a single Docker image. Whether you're a DevOps engineer, an SRE, or a developer, K8sToolbox provides you with everything you need to streamline your workflow and tackle issues in your Kubernetes environment with ease.

## Key Features
- **Network Debugging Tools**: Includes common utilities like `ping`, `curl`, `traceroute`, `netcat`, and more for troubleshooting connectivity issues between Kubernetes pods and services.
- **Pod and Cluster Management**: Tools like `kubectl`, `stern`, and `k9s` make it easier to interact with and manage Kubernetes resources.
- **Database Clients**: Access to clients like `psql`, `mysql`, and `redis-cli` for managing databases within the cluster.
- **Monitoring and Troubleshooting**: Tools like `tcpdump`, `htop`, `iftop`, and `promtool` for advanced monitoring, debugging, and metrics validation.

## Quick Start

### Running Locally
To get started with K8sToolbox locally, you can use Docker:

```sh
docker run -it narmidm/k8stoolbox:latest
```

### Deploying in Kubernetes
To deploy K8sToolbox as a debug pod in your Kubernetes cluster:

```sh
kubectl apply -f manifests/debug-pod.yaml
```

Or use the Helm chart for more flexible deployment:

```sh
helm repo add k8stoolbox https://github.com/narmidm/K8sToolbox/charts
helm install k8stoolbox k8stoolbox/k8stoolbox
```

## Tools Included
- **Network Debugging**: `curl`, `ping`, `traceroute`, `netcat`, `tcpdump`, `iperf3`, `nslookup`, `dig`
- **Pod and Cluster Management**: `kubectl`, `stern`, `k9s`, `jq`, `kubens`, `kubectx`
- **Advanced Debugging**: `strace`, `htop`, `iftop`, `netstat`
- **Storage Management**: `rsync`, `mc` (MinIO Client)
- **Database Clients**: `psql`, `mysql`, `redis-cli`, `mongo`
- **TLS and Security**: `openssl`, `gpg`
- **Editors**: `vim`, `nano`
- **Monitoring**: `promtool`, `curl` for Prometheus/Grafana

## Contributing
We welcome contributions from the community! Feel free to open an issue or submit a pull request if you'd like to add new features or report a bug. Check out our [CONTRIBUTING.md](CONTRIBUTING.md) for more details on how to get involved.

## License
K8sToolbox is licensed under the [Apache License 2.0](LICENSE).

## Inspiration
K8sToolbox was inspired by the [swiss-army-knife](https://github.com/leodotcloud/swiss-army-knife) repository, which serves as a useful multi-purpose tool for DevOps. Our goal is to build upon that foundation and create a specialized, Kubernetes-focused toolkit that helps users effectively troubleshoot and manage their clusters.
