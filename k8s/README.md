# Kubernetes Deployment Guide

## Table of Contents

- [Prerequisites](#prerequisites)
  - [Install Docker](#install-docker)
  - [Install Minikube](#install-minikube)
  - [Install kubectl](#install-kubectl)
- [Application Deployment](#application-deployment)
  - [1. Build Images](#1-build-images)
  - [2. Load Images to Minikube](#2-load-images-to-minikube)
  - [3. Create Namespace](#3-create-namespace)
  - [4. Deploy Services](#4-deploy-services)
  - [5. Configure Ingress](#5-configure-ingress)
  - [6. Monitoring and Debugging](#6-monitoring-and-debugging)
  - [7. Update Images](#7-update-images)
- [Monitoring with k9s](#monitoring-with-k9s)
  - [Install k9s](#install-k9s)
  - [Common Operations](#common-operations)
  - [Common Debugging Flow](#common-debugging-flow)
  - [Tips](#tips)
- [Testing Services](#testing-services)
  - [1. Check Services Status](#1-check-services-status)
  - [2. Test Individual Services](#2-test-individual-services)
  - [3. Test Through Ingress](#3-test-through-ingress)
  - [4. View Service Logs](#4-view-service-logs)
  - [5. Debug Services](#5-debug-services)
- [Useful Links](#useful-links)

## Quick Commands

For common operations, you can use these make commands:

```bash
# Start everything
make all

# Check status
make status

# Restart individual services after code changes
make restart-auth      # Rebuild and restart auth service
make restart-user      # Rebuild and restart user service
make restart-chat      # Rebuild and restart chat service
make restart-gateway   # Rebuild and restart gateway service
make restart-notification  # Rebuild and restart notification service
make restart-subscriber   # Rebuild and restart subscriber service

# Stop everything
make stop

# Clean up everything
make clean
```

[Back to top](#table-of-contents)

### Install Docker

#### MacOS

```bash
# Install Docker Desktop
brew install --cask docker

# Start Docker Desktop
open -a Docker

# Wait until Docker Desktop is fully started
#    The whale icon in menu bar should stop spinning
#    This may take a minute or two

# Set correct Docker context
docker context use desktop-linux

# Verify Docker is accessible
docker ps
```

#### Linux

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Start Docker daemon
sudo systemctl start docker

# Add your user to docker group (optional, to run docker without sudo)
sudo usermod -aG docker $USER
newgrp docker

# Verify Docker is running
docker ps
```

### Install Minikube

1. Install the binary:

```bash
# MacOS
brew install minikube

# Linux
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube
```

2. Ensure Docker is running and context is set correctly:

```bash
# MacOS:
# 1. Open Docker Desktop application
open -a Docker

# 2. Wait until Docker Desktop is fully started
#    The whale icon in menu bar should stop spinning
#    This may take a minute or two

# 3. Set correct Docker context
docker context use desktop-linux

# 4. Verify Docker is accessible
#    If this fails, wait longer for Docker Desktop to fully start
docker ps

# Linux:
sudo systemctl status docker
docker context use default
docker ps
```

3. Start Minikube:

```bash
# Start minikube with Docker driver
minikube start --driver=docker

# After minikube starts successfully, connect to its Docker daemon:
eval $(minikube docker-env)
```

4. Verify everything is running:

```bash
# Check minikube status
minikube status

# Check you can access minikube's Docker
docker ps
```

### Install kubectl

1. https://kubernetes.io/ru/docs/tasks/tools/install-kubectl/
2. _OPTIONAL_: [quick kubernetes namespace and context switcher](https://github.com/blendle/kns)

[Back to top](#table-of-contents)

## Application Deployment

### Prerequisites

> **Important Note**: Before running any commands, ensure you're using the correct Docker context. All commands in this guide will work regardless of your Docker engine (Docker Desktop or native Linux Docker), but the correct context must be set first:
>
> - For MacOS: `docker context use desktop-linux`
> - For Linux: `docker context use default`

### 1. Build Images

```bash
# Connect to minikube's Docker daemon
eval $(minikube docker-env)

# Build images for all services
make docker-build

# Verify built images
docker images | grep msa-messenger
```

### 2. Load Images to Minikube

```bash
# Add images to minikube registry
minikube image load msa-messenger/auth:latest
minikube image load msa-messenger/chat:latest
minikube image load msa-messenger/gateway:latest
minikube image load msa-messenger/notification:latest
minikube image load msa-messenger/subscriber:latest
minikube image load msa-messenger/user:latest

# Verify loaded images
minikube image ls --format table
```

### 3. Create Namespace

```bash
kubectl create namespace messenger
kubectl config set-context --current --namespace=messenger
```

### 4. Deploy Services

Use `make deploy` to deploy all services or deploy services manually:

```bash
# Create services (ClusterIP)
kubectl apply -f k8s/auth/service_cluster_ip.yaml
kubectl apply -f k8s/chat/service_cluster_ip.yaml
kubectl apply -f k8s/gateway/service_cluster_ip.yaml
kubectl apply -f k8s/notification/service_cluster_ip.yaml
kubectl apply -f k8s/subscriber/service_cluster_ip.yaml
kubectl apply -f k8s/user/service_cluster_ip.yaml

# Verify services
kubectl get svc

# Create deployments
kubectl apply -f k8s/auth/deployment.yaml
kubectl apply -f k8s/chat/deployment.yaml
kubectl apply -f k8s/gateway/deployment.yaml
kubectl apply -f k8s/notification/deployment.yaml
kubectl apply -f k8s/subscriber/deployment.yaml
kubectl apply -f k8s/user/deployment.yaml

# Verify deployments and pods
kubectl get deployments
kubectl get pods --show-labels
```

### 5. Configure Ingress

```bash
# Enable Ingress addon in minikube
minikube addons enable ingress

# Verify Ingress controller is running
kubectl get pods -n ingress-nginx

# Create Ingress resource
kubectl apply -f k8s/ingress.yaml

# Verify Ingress
kubectl get ingress

# Add DNS record to /etc/hosts
echo "$(minikube ip) messenger.local" | sudo tee -a /etc/hosts

# Test API availability through Ingress
# Option 1:
curl -H "Host: messenger.local" http://$(minikube ip)/api/v1/health

# Option 2 (if option 1 doesn't work in your shell):
MINIKUBE_IP=$(minikube ip)
curl -H "Host: messenger.local" http://$MINIKUBE_IP/api/v1/health
```

### 6. Monitoring and Debugging

```bash
# View pod logs
kubectl logs -f deployment/gateway

# Describe resources
kubectl describe deployment gateway
kubectl describe service gateway
kubectl describe ingress

# Port forwarding for debugging
kubectl port-forward service/gateway 8080:8080
```

### 7. Update Images

```bash
# Update image in deployment
kubectl set image deployment/gateway gateway-container=msa-messenger/gateway:latest
```

[Back to top](#table-of-contents)

## Monitoring with k9s

### Install k9s

```bash
# Install via Homebrew
brew install k9s

# Start k9s
k9s
```

### Common Operations

#### Basic Navigation

```bash
# Switch context
:context

# Switch namespace
:namespace
# or just press Ctrl+A to show all namespaces

# Filter resources
/ + type filter text
```

#### View Resources

```bash
# View all pods
:pods or :po

# View deployments
:deploy or :dp

# View services
:svc

# View logs
l (when pod is selected)

# Describe resource
d (when resource is selected)
```

#### Pod Management

```bash
# Shell into container
s (when pod is selected)

# View pod logs
l (when pod is selected)

# Restart pod
ctrl + k (when pod is selected)

# Scale deployment
s (when deployment is selected)
```

#### Useful Shortcuts

```bash
# Help
?

# Quit
:q or ctrl + c

# Previous page
esc

# Port forward
shift + f (when service/pod is selected)
```

### Common Debugging Flow

1. Start k9s: `k9s`
2. Switch to pods view: `:pods`
3. Check pod status and ready state
4. View logs if pod is not ready: `l`
5. Describe pod for more details: `d`
6. Shell into pod if needed: `s`
7. Restart pod if required: `ctrl + k`

### Tips

- Use port forwarding to test individual services
- Watch pod logs in real-time for debugging
- Use label filters to find specific resources
- Check events when pods fail to start
- Monitor resource usage with `:pulses`

[Back to top](#table-of-contents)

## Testing Services

### 1. Check Services Status

```bash
# Check all pods are running
kubectl get pods -n messenger

# Check all services are created
kubectl get svc -n messenger

# Check ingress is configured
kubectl get ingress -n messenger
```

### 2. Test Individual Services

You can test each service directly using port-forward. Note: Run port-forward and curl commands in separate terminals, as port-forward blocks the terminal while running.

```bash
# Auth Service
# Terminal 1:
kubectl port-forward service/auth-service 8081:8081 -n messenger
# Terminal 2:
curl http://localhost:8081/health

# User Service
# Terminal 1:
kubectl port-forward service/user-service 8082:8082 -n messenger
# Terminal 2:
curl http://localhost:8082/health

# Chat Service
# Terminal 1:
kubectl port-forward service/chat-service 8083:8083 -n messenger
# Terminal 2:
curl http://localhost:8083/health

# Subscriber Service
# Terminal 1:
kubectl port-forward service/subscriber-service 8084:8084 -n messenger
# Terminal 2:
curl http://localhost:8084/health

# Notification Service
# Terminal 1:
kubectl port-forward service/notification-service 8085:8085 -n messenger
# Terminal 2:
curl http://localhost:8085/health

# Gateway Service
# Terminal 1:
kubectl port-forward service/gateway-service 8080:8080 -n messenger
# Terminal 2:
curl http://localhost:8080/health
```

### 3. Test Through Ingress

```bash
# Test API availability through Ingress
curl -H "Host: messenger.local" http://$(minikube ip)/api/v1/health
```

### 4. View Service Logs

```bash
# View pod logs
kubectl logs -f deployment/gateway
```

### 5. Debug Services

```bash
# Debug services
kubectl port-forward service/gateway 8080:8080 -n messenger
```

[Back to top](#table-of-contents)

## Useful Links
