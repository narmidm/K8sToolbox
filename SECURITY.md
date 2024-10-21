# Security Policy for K8sToolbox

## Reporting a Vulnerability

If you believe you've found a security vulnerability in K8sToolbox, we encourage you to responsibly disclose it by following the procedure below. Your efforts in making our project more secure are highly appreciated.

### How to Report

1. **Email Contact**: Send an email to our security team at [imranaec@outlook.com](mailto:imranaec@outlook.com) with a detailed description of the vulnerability, including:
    - Steps to reproduce the issue.
    - Affected versions of K8sToolbox.
    - Any potential impacts or exploitation risks.
    - Your suggested fix or mitigation, if any.

2. **Confidentiality**: Please keep the details of any potential security issue confidential until we have addressed it.

3. **Acknowledgement**: We will acknowledge your report within 3 business days and work with you to fully understand the issue and its implications.

4. **Fix and Disclosure**: Once the vulnerability is verified, we will work on a patch and schedule a release. We will coordinate with you for responsible disclosure and mention your contribution unless you prefer to remain anonymous.

## Supported Versions

We provide security updates for the following versions of K8sToolbox:

| Version        | Supported         |
|----------------|-------------------|
| Latest (master)| Yes               |
| v1.0 and above | Yes               |

Please upgrade to a supported version to receive the latest security fixes.

## Security Best Practices

To ensure the security of your Kubernetes environment when using K8sToolbox, please follow these best practices:

- **Run K8sToolbox in a Secure Namespace**: Isolate the tool in a secure namespace with appropriate Role-Based Access Control (RBAC) rules to limit its permissions.

- **Limit Privileged Operations**: Only run K8sToolbox with the necessary privileges. Avoid running it in `privileged` mode unless absolutely required for debugging purposes.

- **Audit Permissions**: Regularly audit the permissions granted to the `k8stoolbox` ServiceAccount and ClusterRoleBinding. Ensure that they adhere to the principle of least privilege.

- **Use Trusted Networks**: Deploy K8sToolbox in a trusted network and restrict access to the Pod via Kubernetes Network Policies to reduce the attack surface.

- **Keep Up-to-Date**: Always use the latest version of K8sToolbox to benefit from security patches and updates.

## Security Response Team

Our security response team is dedicated to resolving security issues and ensuring a safe experience for all users. Please do not publicly disclose any security-related information until we provide a solution.

For general security inquiries, contact us at [imranaec@outlook.com](mailto:imranaec@outlook.com).

## Code Security and Hardening Measures

We are committed to maintaining the highest levels of security within the K8sToolbox codebase. Below are some key measures we've taken to ensure the tool's security:

- **Code Reviews**: All changes are subject to thorough code reviews, with an emphasis on security and stability.
- **Dependency Management**: We regularly audit and update third-party libraries to prevent vulnerabilities stemming from outdated dependencies.
- **Static Analysis**: We use automated tools for static code analysis to identify potential vulnerabilities during development.

Thank you for helping to improve the security of K8sToolbox.
