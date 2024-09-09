##@ Check Environment Variables
.PHONY: check-env
check-env: ## Check environment variables
	@echo "Checking environment variables..."
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found. Please create one with required variables."; \
		exit 1; \
	fi
	@echo "Loading .env file..."
	@set -a; . ./.env; set +a
	@echo "All required environment variables are set."

# Include and export variables from .env
-include .env
export $(shell sed 's/=.*//' .env)

# Set variables
CONTROLLER_IMG ?= quay.io/climatik-project/climatik-controller
WEBHOOK_IMG ?= quay.io/climatik-project/webhook

CLUSTER_PROVIDER ?= kind
LOCAL_DEV_CLUSTER_VERSION ?= main
KIND_WORKER_NODES ?= 2

# Go related variables
GOENV = GO111MODULE="" \
        GOOS=$(shell go env GOOS) \
        GOARCH=$(shell go env GOARCH)

# Output directory
OUTPUT_DIR := _output
CROSS_BUILD_BINDIR := $(OUTPUT_DIR)/bin

.DEFAULT_GOAL := help

export CTR_CMD     ?= $(or $(shell command -v podman), $(shell command -v docker))

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.DEFAULT_GOAL := default

.PHONY: default
default: check-env tests build-controller-image push-controller-image build-webhook-image push-webhook-image
.PHONY: tests
tests: ## Run tests
	PROMETHEUS_HOST="http://127.0.0.1:9090" python -m unittest discover python/tests
	go test -v ./internal/alert/tests

.PHONY: build
build: ## Build the project
	$(GOENV) go build -o bin/manager ./cmd/controller/main.go
	$(GOENV) go build -o bin/webhook ./cmd/webhook/main.go

.PHONY: run
run: ## Run the project
	$(GOENV) go run ./cmd/...

.PHONY: local-run
local-run: ## Run the project locally
	kopf run python/climatik_operator/operator.py --verbose

##@ Build
.PHONY: build-controller-image
build-controller-image: tests ## Build Docker image for Container Registry
ifeq ($(OPT), NO_CACHE_BUILD)
	$(CTR_CMD) build --no-cache -f dockerfiles/Dockerfile -t $(CONTROLLER_IMG):latest .
else
	$(CTR_CMD) build -f dockerfiles/Dockerfile -t $(CONTROLLER_IMG):latest .
endif

.PHONY: build-webhook-image
build-webhook-image: ## Build Docker image for the webhook
ifeq ($(OPT), NO_CACHE_BUILD)
	$(CTR_CMD) build --no-cache -f dockerfiles/Dockerfile.webhook -t $(WEBHOOK_IMG):latest .
else
	$(CTR_CMD) build -f dockerfiles/Dockerfile.webhook -t $(WEBHOOK_IMG):latest .
endif

##@ Push Image

.PHONY: push-controller-image
push-controller-image: build-controller-image ## Push Docker image to Container Registry
	$(CTR_CMD) push $(CONTROLLER_IMG):latest

.PHONY: push-webhook-image
push-webhook-image: build-webhook-image ## Push Docker image for the webhook to Container Registry
	$(CTR_CMD) push $(WEBHOOK_IMG):latest

##@ Clean Up Resources

.PHONY: clean-up
clean-up: ## Clean up deployments
	kubectl delete deployment operator-powercapping-controller-manager -n operator-powercapping-system --ignore-not-found
	kubectl delete deployment operator-powercapping-webhook-manager -n operator-powercapping-system --ignore-not-found
	kubectl delete deployment llama2-7b -n operator-powercapping-system --ignore-not-found
	kubectl delete deployment mistral-7b -n operator-powercapping-system --ignore-not-found
	kubectl delete deployment stress -n operator-powercapping-system --ignore-not-found
	if kubectl get crd scaledobjects.keda.sh > /dev/null 2>&1; then \
		kubectl delete scaledobject mistral-7b-scaleobject -n operator-powercapping-system --ignore-not-found; \
		kubectl delete scaledobject llama2-7b-scaleobject -n operator-powercapping-system --ignore-not-found; \
		kubectl delete scaledobject stress-scaleobject -n operator-powercapping-system --ignore-not-found; \
	else \
		echo "ScaledObject resource type not found, ignoring."; \
	fi

##@ Secrets

.PHONY: create-env-secret
create-env-secret: ## Create env secret from .env file
	kubectl create namespace operator-powercapping-system --dry-run=client -o yaml | kubectl apply -f -
	kubectl create namespace system --dry-run=client -o yaml | kubectl apply -f -
	kubectl delete secret env-secrets -n operator-powercapping-system --ignore-not-found
	kubectl delete secret env-secrets -n system --ignore-not-found

	kubectl create secret generic env-secrets --from-env-file=.env -n system
	kubectl create secret generic env-secrets --from-env-file=.env -n operator-powercapping-system
	@echo "Environment secret has been recreated in the operator-powercapping-system namespace."
	@echo "Environment secret has been recreated in the system namespace."


##@ Deployment

.PHONY: deploy-config
deploy-config: deploy ## Deploy with modified config

.PHONY: deploy-config
deploy-config: deploy ## Deploy with modified config using Container Registry

.PHONY: deploy
deploy: create-env-secret ## Deploy to Kubernetes
	set -x
	kubectl apply -f config/crd/bases
	kustomize build config/default | kubectl apply -f -
	kubectl apply -f deploy/climatik-operator/manifests/crd.yaml
	kubectl apply -f deploy/climatik-operator/manifests/sample-powercappingconfig.yaml
	kubectl apply -f deploy/climatik-operator/manifests/manager-deployment.yaml
	kubectl apply -f deploy/climatik-operator/manifests/webhook-deployment.yaml

.PHONY: deploy-sample
deploy-sample: clean-up ## Deploy to Kubernetes using Container Registry
	kubectl apply -f config/crd/bases
	kustomize build config/default | kubectl apply -f -
	kubectl apply -f deploy/climatik-operator/manifests/crd.yaml
	kubectl apply -f hack/keda/keda-2.10.0.yaml
	kubectl wait --for=condition=Available --timeout=600s apiservice v1beta1.external.metrics.k8s.io
	kubectl apply -f deploy/climatik-operator/manifests/deployment-mistral-7b.yaml
	kubectl apply -f deploy/climatik-operator/manifests/deployment-llama2-7b.yaml
	kubectl apply -f deploy/climatik-operator/manifests/deployment-stress.yaml
	kubectl apply -f deploy/climatik-operator/manifests/scaleobject.yaml
	kubectl apply -f deploy/climatik-operator/manifests/sample-powercappingconfig.yaml
	kubectl apply -f deploy/climatik-operator/manifests/manager-deployment.yaml
	kubectl apply -f deploy/climatik-operator/manifests/webhook-deployment.yaml

##@ Start Cluster

.PHONY: cluster-up
cluster-up: ## Setup a cluster for local development
	CLUSTER_PROVIDER=$(CLUSTER_PROVIDER) \
	LOCAL_DEV_CLUSTER_VERSION=$(LOCAL_DEV_CLUSTER_VERSION) \
	KIND_WORKER_NODES=$(KIND_WORKER_NODES) \
	BUILD_CONTAINERIZED=$(BUILD_CONTAINERIZED) \
	PROMETHEUS_ENABLE=true \
	GRAFANA_ENABLE=true \
	./hack/cluster.sh up
