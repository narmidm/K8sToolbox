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
FROM alpine:latest

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
    bash

# Install kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && rm kubectl

# Install stern
RUN curl -L -o /usr/local/bin/stern https://github.com/stern/stern/releases/download/v1.22.0/stern_linux_amd64 && \
    chmod +x /usr/local/bin/stern

# Install k9s
RUN curl -L -o /usr/local/bin/k9s https://github.com/derailed/k9s/releases/download/v0.27.3/k9s_Linux_x86_64.tar.gz && \
    tar -xz -C /usr/local/bin/ -f /usr/local/bin/k9s && chmod +x /usr/local/bin/k9s

# Install kubens and kubectx
RUN git clone https://github.com/ahmetb/kubectx /opt/kubectx && \
    ln -s /opt/kubectx/kubens /usr/local/bin/kubens && \
    ln -s /opt/kubectx/kubectx /usr/local/bin/kubectx

# Install mc (MinIO Client)
RUN curl -O https://dl.min.io/client/mc/release/linux-amd64/mc && \
    chmod +x mc && \
    mv mc /usr/local/bin/

# Copy the Golang executable from the builder stage
COPY --from=builder /app/k8stoolbox /usr/local/bin/k8stoolbox

# Clean up unnecessary files to keep the image small
RUN rm -rf /var/cache/apk/*

# Set the default command to an interactive shell for general debugging and troubleshooting
CMD ["/bin/sh"]
