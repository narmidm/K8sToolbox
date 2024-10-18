# K8sToolbox

**Your Swiss Army Knife for Kubernetes Troubleshooting**

K8sToolbox is a versatile collection of tools for debugging, troubleshooting, and managing Kubernetes clusters, all bundled into a single Docker image. Whether you're a DevOps engineer, an SRE, or a developer, K8sToolbox provides you with everything you need to streamline your workflow and tackle issues in your Kubernetes environment with ease.

## Key Features
- **Network Debugging Tools**: Includes common utilities like `ping`, `curl`, `traceroute`, and `nslookup` for troubleshooting connectivity issues between Kubernetes pods and services.
- **Pod Inspection and Interaction**: Quickly inspect resources, check pod logs, execute commands, and troubleshoot Kubernetes objects more efficiently.
- **Deployment Flexibility**: Deploy K8sToolbox as a standalone pod, as a DaemonSet across nodes, or even as a sidecar for deeper insights into specific applications.
- **Ease of Use**: Easily run diagnostics within your cluster without the need to install multiple tools.

## Quick Start

### Running Locally
To get started with K8sToolbox locally, you can use Docker:

```sh
docker run -it k8stoolbox:latest
```

### Deploying in Kubernetes
To deploy K8sToolbox as a debug pod in your Kubernetes cluster:

```sh
kubectl apply -f https://raw.githubusercontent.com/narmidm/K8sToolbox/main/k8s/debug-pod.yaml
```

## Inspiration
K8sToolbox was inspired by the [swiss-army-knife](https://github.com/leodotcloud/swiss-army-knife) repository, which serves as a useful multi-purpose tool for DevOps. Our goal is to build upon that foundation and create a specialized, Kubernetes-focused toolkit that helps users effectively troubleshoot and manage their clusters.

## Contributing
We welcome contributions from the community! Feel free to open an issue or submit a pull request if you'd like to add new features or report a bug. Check out our [CONTRIBUTING.md](CONTRIBUTING.md) for more details on how to get involved.

## License
K8sToolbox is licensed under the [MIT License](LICENSE).
