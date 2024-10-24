# Dockerfile for K8sToolbox

# Stage 1: Builder - Build the Golang tool
FROM golang:1.22.0-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git build-base

# Set the working directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . ./

# Build the Golang executable
RUN go build -o k8stoolbox main.go

# Stage 2: Final Image
FROM alpine:3.20.3

# Set environment variables
ENV LANG=C.UTF-8

# Install essential tools
RUN apk add --no-cache \
    curl \
    iproute2 \
    iputils \
    netcat-openbsd \
    tcpdump \
    bind-tools \
    traceroute \
    iperf3 \
    jq \
    strace \
    htop \
    iftop \
    net-tools \
    rsync \
    openssl \
    gpg \
    vim \
    nano \
    busybox-extras \
    mariadb-client \
    postgresql-client \
    redis \
    mongodb-tools \
    helm \
    socat \
    ncdu \
    bash \
    ca-certificates \
    conntrack-tools \
    ethtool \
    iptables \
    less \
    mtr \
    openssh-client \
    psmisc \
    tcptraceroute \
    ngrep \
    yq \
    nmap \
    screen \
    tmux

# Install kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && rm kubectl

# Install stern
RUN curl -L -o /usr/local/bin/stern https://github.com/stern/stern/releases/download/v1.22.0/stern_linux_amd64 && \
    chmod +x /usr/local/bin/stern

# Install k9s (Fixed Download and Extraction)
RUN curl -L -o /tmp/k9s_Linux_amd64.tar.gz https://github.com/derailed/k9s/releases/download/v0.32.5/k9s_Linux_amd64.tar.gz && \
    tar -xz -C /usr/local/bin/ -f /tmp/k9s_Linux_amd64.tar.gz k9s && chmod +x /usr/local/bin/k9s && rm /tmp/k9s_Linux_amd64.tar.gz

# Install mc (MinIO Client)
RUN curl -O https://dl.min.io/client/mc/release/linux-amd64/mc && \
    chmod +x mc && \
    mv mc /usr/local/bin/

# Set the working directory
WORKDIR /root/

# Copy the Golang executable from the builder stage
COPY --from=builder /app/k8stoolbox /usr/local/bin/k8stoolbox

# Copy helper scripts
COPY scripts/ /usr/local/bin/scripts/

# Ensure scripts are executable
RUN chmod +x /usr/local/bin/scripts/*.sh

# Create symlinks for easy access to scripts
RUN ln -s /usr/local/bin/scripts/aggregate_logs.sh /usr/local/bin/aggregate_logs && \
    ln -s /usr/local/bin/scripts/auto_recover.sh /usr/local/bin/auto_recover && \
    ln -s /usr/local/bin/scripts/auto_scaling.sh /usr/local/bin/auto_scaling && \
    ln -s /usr/local/bin/scripts/backup_restore.sh /usr/local/bin/backup_restore && \
    ln -s /usr/local/bin/scripts/clean_stale_resources.sh /usr/local/bin/clean_stale_resources && \
    ln -s /usr/local/bin/scripts/connectivity_test.sh /usr/local/bin/connectivity_test && \
    ln -s /usr/local/bin/scripts/delete_stuck_crds.sh /usr/local/bin/delete_stuck_crds && \
    ln -s /usr/local/bin/scripts/delete_stuck_namespace.sh /usr/local/bin/delete_stuck_namespace && \
    ln -s /usr/local/bin/scripts/healthcheck.sh /usr/local/bin/healthcheck && \
    ln -s /usr/local/bin/scripts/network_diag.sh /usr/local/bin/network_diag && \
    ln -s /usr/local/bin/scripts/resource_usage.sh /usr/local/bin/resource_usage && \
    ln -s /usr/local/bin/scripts/restart_failed_pods.sh /usr/local/bin/restart_failed_pods && \
    ln -s /usr/local/bin/scripts/snapshot_audit.sh /usr/local/bin/snapshot_audit && \
    ln -s /usr/local/bin/scripts/test_network_policy.sh /usr/local/bin/test_network_policy

# Clean up unnecessary files to keep the image small
RUN rm -rf /var/cache/apk/*

# Run the Go application
CMD ["tail", "-f", "/dev/null"]