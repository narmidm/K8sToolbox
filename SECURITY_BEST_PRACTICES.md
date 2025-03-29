# K8sToolbox Security Best Practices

This document outlines security best practices for deploying and using K8sToolbox in your Kubernetes environment.

## Table of Contents
- [Introduction](#introduction)
- [Deployment Security](#deployment-security)
  - [RBAC Configuration](#rbac-configuration)
  - [Container Security Context](#container-security-context)
  - [Host Access](#host-access)
- [Runtime Security](#runtime-security)
  - [Script Execution](#script-execution)
  - [Network Access](#network-access)
  - [Authentication](#authentication)
- [Production Recommendations](#production-recommendations)
- [Security Checklist](#security-checklist)

## Introduction

K8sToolbox is a powerful utility for Kubernetes management and troubleshooting. With great power comes great responsibility - many of the capabilities that make it useful for debugging can also create security risks if not properly managed.

## Deployment Security

### RBAC Configuration

K8sToolbox includes two RBAC configurations:

1. **Restricted Permissions (Recommended)**:
   - Uses the principle of least privilege
   - Grants only the permissions needed for K8sToolbox to function
   - Recommended for all environments, especially production

2. **Cluster Admin (Not Recommended)**:
   - Grants full administrative access to the cluster
   - Only use in isolated, non-production environments
   - Creates significant security risk

**Best Practice:** Always use restricted permissions by setting `security.useRestrictedPermissions=true` in Helm values.

```yaml
security:
  useRestrictedPermissions: true
```

### Container Security Context

The security context defines privilege and access controls for the container.

**Best Practice:** Configure the security context with the following settings:

```yaml
security:
  podSecurityContext:
    fsGroup: 10001
    runAsUser: 10001
    runAsGroup: 10001
    runAsNonRoot: true
  
  containerSecurityContext:
    privileged: false
    allowPrivilegeEscalation: false
    readOnlyRootFilesystem: true
    capabilities:
      drop:
        - ALL
      add:
        - NET_ADMIN  # Only if needed for network diagnostics
        - NET_RAW    # Only if needed for network diagnostics
```

**Note:** Some network diagnostic tools may require `NET_ADMIN` and `NET_RAW` capabilities. If these aren't needed, don't add them.

### Host Access

K8sToolbox can be configured to access the host filesystem.

**Best Practice:** Avoid mounting the host filesystem unless absolutely necessary.

```yaml
volumes:
  mountHostRoot: false
```

If host access is required, consider:
1. Mounting only specific directories instead of the entire root filesystem
2. Making the mounts read-only where possible
3. Using this feature only in non-production environments

## Runtime Security

### Script Execution

K8sToolbox provides scripts that can modify cluster resources.

**Best Practices:**
1. Only run scripts in namespaces where changes are intended
2. Consider executing scripts with `--dry-run` first to preview changes
3. Implement proper access controls to limit who can execute these scripts
4. Always review script output to understand what changed

### Network Access

Some K8sToolbox scripts provide network diagnostic capabilities that can potentially be used for unintended purposes.

**Best Practices:**
1. Use Network Policies to limit which pods K8sToolbox can communicate with
2. Only deploy K8sToolbox in namespaces where it's needed
3. Restrict egress traffic from the K8sToolbox pod when in production

Example Network Policy:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: k8stoolbox-network-policy
spec:
  podSelector:
    matchLabels:
      app: k8stoolbox
  policyTypes:
  - Egress
  egress:
  - to:
    - namespaceSelector:
        matchLabels:
          purpose: production
```

### Authentication

When using the K8sToolbox web interface (if enabled), proper authentication is essential.

**Best Practices:**
1. Always enable authentication when using the web interface
2. Generate strong passwords and rotate them regularly
3. Consider integrating with existing identity providers
4. Use HTTPS for all traffic to the web interface

```yaml
auth:
  enabled: true
  generatePassword: true
```

## Production Recommendations

For production environments, follow these additional guidelines:

1. **Namespace Isolation**: Deploy K8sToolbox in its own namespace with tight access controls
2. **Resource Limits**: Set appropriate CPU and memory limits to prevent resource exhaustion
3. **Audit Logging**: Enable Kubernetes audit logging to track actions performed by K8sToolbox
4. **Regular Updates**: Keep K8sToolbox updated to the latest version to benefit from security patches
5. **Scheduled Deployments**: Consider deploying K8sToolbox only when needed and removing it afterward
6. **Image Scanning**: Scan the K8sToolbox image for vulnerabilities before deployment

## Security Checklist

Use this checklist before deploying K8sToolbox:

- [ ] RBAC is set to restricted permissions
- [ ] Container is not running as privileged
- [ ] Host root filesystem is not mounted
- [ ] Non-root user is configured
- [ ] Resource limits are defined
- [ ] Network policies are implemented
- [ ] Authentication is enabled for the web interface
- [ ] Image has been scanned for vulnerabilities
- [ ] Read-only root filesystem is enabled
- [ ] Unnecessary capabilities are dropped
- [ ] Deployment is in an isolated namespace

By following these security best practices, you can ensure that K8sToolbox enhances your Kubernetes management capabilities without introducing unnecessary risk to your environment. 