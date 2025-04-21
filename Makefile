# Before running any commands, ensure you're using the correct Docker context:
# For MacOS:
# 1. Open Docker Desktop application
# 2. Set correct Docker context in terminal	:
#    docker context use desktop-linux
# For Linux:
# Set correct Docker context in terminal:
# docker context use default

.PHONY: start stop clean build docker-build deploy all restart-auth restart-user restart-chat restart-gateway restart-notification restart-subscriber

# Build services with Go
build:
	go build -o bin/auth ./auth/cmd
	go build -o bin/chat ./chat/cmd
	go build -o bin/gateway ./gateway/cmd
	go build -o bin/notification ./notification/cmd
	go build -o bin/subscriber ./subscriber/cmd
	go build -o bin/user ./user/cmd

# Build Docker images
docker-build:
	docker build -t msa-messenger/auth:latest -f auth/Dockerfile .
	docker build -t msa-messenger/chat:latest -f chat/Dockerfile .
	docker build -t msa-messenger/gateway:latest -f gateway/Dockerfile .
	docker build -t msa-messenger/notification:latest -f notification/Dockerfile .
	docker build -t msa-messenger/subscriber:latest -f subscriber/Dockerfile .
	docker build -t msa-messenger/user:latest -f user/Dockerfile .

# Rebuild and restart individual services
restart-auth:
	docker build -t msa-messenger/auth:latest -f auth/Dockerfile .
	kubectl rollout restart deployment auth-service -n messenger
	@echo "Waiting for auth service to restart..."
	kubectl rollout status deployment auth-service -n messenger

restart-user:
	docker build -t msa-messenger/user:latest -f user/Dockerfile .
	kubectl rollout restart deployment user-service -n messenger
	@echo "Waiting for user service to restart..."
	kubectl rollout status deployment user-service -n messenger

restart-chat:
	docker build -t msa-messenger/chat:latest -f chat/Dockerfile .
	kubectl rollout restart deployment chat-service -n messenger
	@echo "Waiting for chat service to restart..."
	kubectl rollout status deployment chat-service -n messenger

restart-gateway:
	docker build -t msa-messenger/gateway:latest -f gateway/Dockerfile .
	kubectl rollout restart deployment gateway-deployment -n messenger
	@echo "Waiting for gateway service to restart..."
	kubectl rollout status deployment gateway-deployment -n messenger

restart-notification:
	docker build -t msa-messenger/notification:latest -f notification/Dockerfile .
	kubectl rollout restart deployment notification-service -n messenger
	@echo "Waiting for notification service to restart..."
	kubectl rollout status deployment notification-service -n messenger

restart-subscriber:
	docker build -t msa-messenger/subscriber:latest -f subscriber/Dockerfile .
	kubectl rollout restart deployment subscriber-service -n messenger
	@echo "Waiting for subscriber service to restart..."
	kubectl rollout status deployment subscriber-service -n messenger

# Start minikube
start:
	minikube start --driver=docker

# Deploy all services
deploy:
	kubectl create namespace messenger || true
	kubectl config set-context --current --namespace=messenger
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/auth/service_cluster_ip.yaml
	kubectl apply -f k8s/chat/service_cluster_ip.yaml
	kubectl apply -f k8s/gateway/service_cluster_ip.yaml
	kubectl apply -f k8s/notification/service_cluster_ip.yaml
	kubectl apply -f k8s/subscriber/service_cluster_ip.yaml
	kubectl apply -f k8s/user/service_cluster_ip.yaml
	kubectl apply -f k8s/auth/deployment.yaml
	kubectl apply -f k8s/chat/deployment.yaml
	kubectl apply -f k8s/gateway/deployment.yaml
	kubectl apply -f k8s/notification/deployment.yaml
	kubectl apply -f k8s/subscriber/deployment.yaml
	kubectl apply -f k8s/user/deployment.yaml
	minikube addons enable ingress
	kubectl apply -f k8s/ingress.yaml

# Check status
status:
	kubectl get pods -n messenger
	kubectl get svc -n messenger
	kubectl get ingress -n messenger

# Stop minikube
stop:
	minikube stop

# Clean up everything
clean:
	rm -rf bin/
	kubectl delete namespace messenger || true
	minikube delete

# All-in-one command to start everything (using Docker)
all: start docker-build deploy status

# Show minikube IP and add to /etc/hosts
setup-hosts:
	@echo "Add this line to /etc/hosts:"
	@echo "$$(minikube ip) messenger.local"
	@echo "Run: echo \"$$(minikube ip) messenger.local\" | sudo tee -a /etc/hosts"
