# Set variables
IMG ?= quay.io/climatik-project/climatik-operator
GITHUB_REPO ?= climatik-project
GHCR_IMG ?= ghcr.io/$(GITHUB_USERNAME)/$(GITHUB_REPO)

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
default: tests build-image-ghcr push-image-ghcr modify-manager-yaml clean-up deploy-ghcr ## Run default targets

.PHONY: tests
tests: ## Run tests
	PROMETHEUS_HOST="http://localhost:9090" python -m unittest discover python/tests
	go test -v ./internal/alert/tests

.PHONY: build
build: ## Build the project
	$(GOENV) go build -o bin/manager ./cmd/...

.PHONY: run
run: ## Run the project
	$(GOENV) go run ./cmd/...

.PHONY: local-run
local-run: ## Run the project locally
	kopf run python/climatik_operator/operator.py --verbose

##@ Build

.PHONY: build-image
build-image: tests ## Build Docker image
	$(CTR_CMD) build -t $(IMG):latest .

.PHONY: build-image-ghcr
build-image-ghcr: tests ## Build Docker image for GitHub Container Registry
ifeq ($(OPT), NO_CACHE_BUILD)
	$(CTR_CMD) build --no-cache -t $(GHCR_IMG):latest .
else
	$(CTR_CMD) build -t $(GHCR_IMG):latest .
endif

##@ Deployment

.PHONY: push-image
push-image: build-image ## Push Docker image
	$(CTR_CMD) push $(IMG):latest

.PHONY: push-image-ghcr
push-image-ghcr: build-image-ghcr ## Push Docker image to GitHub Container Registry
	echo $(GITHUB_PAT) | $(CTR_CMD) login ghcr.io -u $(GITHUB_USERNAME) --password-stdin
	$(CTR_CMD) push $(GHCR_IMG):latest

.PHONY: clean-up
clean-up: ## Clean up deployments
	kubectl delete deployment operator-powercapping-controller-manager -n operator-powercapping-system --ignore-not-found
	kubectl delete deployment llama2-7b -n operator-powercapping-system --ignore-not-found
	kubectl delete deployment mistral-7b -n operator-powercapping-system --ignore-not-found
	if kubectl get crd scaledobjects.keda.sh > /dev/null 2>&1; then \
		kubectl delete scaledobject mistral-7b-scaleobject -n operator-powercapping-system --ignore-not-found; \
		kubectl delete scaledobject llama2-7b-scaleobject -n operator-powercapping-system --ignore-not-found; \
	else \
		echo "ScaledObject resource type not found, ignoring."; \
	fi

.PHONY: deploy
deploy: ## Deploy to Kubernetes
	kubectl apply -f config/crd/bases
	kustomize build config/default | kubectl apply -f -
	kubectl apply -f deploy/climatik-operator/manifests/crd.yaml
	kubectl apply -f deploy/climatik-operator/manifests/sample_powercapping.yaml
	file=$$(cat "deploy/climatik-operator/manifests/deployment.yaml" | sed "s/\$${GITHUB_USERNAME}/$(GITHUB_USERNAME)/g" | sed "s/\$${GITHUB_REPO}/$(GITHUB_REPO)/g"); \
	echo "$$file"; \
	echo "$$file" | kubectl apply -f -

.PHONY: deploy-ghcr
deploy-ghcr: clean-up ## Deploy to Kubernetes using GitHub Container Registry
	kubectl apply -f config/crd/bases
	kustomize build config/default | kubectl apply -f -
	kubectl apply -f deploy/climatik-operator/manifests/crd.yaml
	kubectl apply -f hack/keda/keda-2.10.0.yaml
	kubectl wait --for=condition=Available --timeout=600s apiservice v1beta1.external.metrics.k8s.io
	kubectl apply -f deploy/climatik-operator/manifests/deployment-mistral-7b.yaml
	kubectl apply -f deploy/climatik-operator/manifests/deployment-llama2-7b.yaml
	kubectl apply -f deploy/climatik-operator/manifests/scaleobject.yaml
	kubectl apply -f deploy/climatik-operator/manifests/sample_powercapping.yaml
	file=$$(cat "deploy/climatik-operator/manifests/deployment.yaml" | sed "s/\$${GITHUB_USERNAME}/$(GITHUB_USERNAME)/g" | sed "s/\$${GITHUB_REPO}/$(GITHUB_REPO)/g"); \
	echo "$$file"; \
	echo "$$file" | kubectl apply -f -

.PHONY: modify-manager-yaml
modify-manager-yaml: ## Modify manager.yaml
	sed -i.bak "s/\$${GITHUB_USERNAME}/$(GITHUB_USERNAME)/g" config/manager/manager.yaml && rm config/manager/manager.yaml.bak
	sed -i.bak "s/\$${GITHUB_REPO}/$(GITHUB_REPO)/g" config/manager/manager.yaml && rm config/manager/manager.yaml.bak

.PHONY: cluster-up
cluster-up: ## Setup a cluster for local development
	CLUSTER_PROVIDER=$(CLUSTER_PROVIDER) \
	LOCAL_DEV_CLUSTER_VERSION=$(LOCAL_DEV_CLUSTER_VERSION) \
	KIND_WORKER_NODES=$(KIND_WORKER_NODES) \
	BUILD_CONTAINERIZED=$(BUILD_CONTAINERIZED) \
	PROMETHEUS_ENABLE=true \
	GRAFANA_ENABLE=true \
	./hack/cluster.sh up

.PHONY: deploy-config
deploy-config: modify-manager-yaml deploy ## Deploy with modified config

.PHONY: deploy-config-ghcr
deploy-config-ghcr: modify-manager-yaml deploy-ghcr ## Deploy with modified config using GitHub Container Registry

##@ Secrets

.PHONY: create-slack-secret
create-slack-secret: ## Create Slack secret from .env file
	kubectl delete secret slack-secrets -n operator-powercapping-system --ignore-not-found
	kubectl create secret generic slack-secrets --from-env-file=.env -n operator-powercapping-system
	@echo "Slack secret has been recreated in the operator-powercapping-system namespace."
