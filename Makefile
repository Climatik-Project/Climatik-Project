.PHONY: all tests build-image-ghcr push-image-ghcr release-ghcr deploy-ghcr cluster-up

# Set variables
IMG ?= quay.io/climatik-project/climatik-operator
GITHUB_USERNAME ?= your-github-username
GITHUB_REPO ?= climatik-project
GHCR_IMG ?= ghcr.io/$(GITHUB_USERNAME)/$(GITHUB_REPO)

CLUSTER_PROVIDER ?= kind
LOCAL_DEV_CLUSTER_VERSION ?= main
KIND_WORKER_NODES ?=2

all: tests build-image-ghcr push-image-ghcr deploy-ghcr

tests:
	cd python && PROMETHEUS_HOST="http://localhost:9090" python -m unittest discover tests

build-image: tests
	docker build -t $(IMG):latest .

build-image-ghcr: tests
	docker build -t $(GHCR_IMG):latest .

push-image: build-image
	docker push $(IMG):latest

push-image-ghcr: build-image-ghcr
	echo $(GITHUB_PAT) | docker login ghcr.io -u $(GITHUB_USERNAME) --password-stdin
	docker push $(GHCR_IMG):latest

release: push-image

release-ghcr: push-image-ghcr

deploy: release
	kubectl apply -f config/crd/bases
	kustomize build config/default | kubectl apply -f -
	kubectl apply -f deploy/climatik-operator/manifests/crd.yaml
	kubectl apply -f deploy/climatik-operator/manifests/sample_powercapping.yaml
	kubectl apply -f deploy/climatik-operator/manifests/deployment.yaml

deploy-ghcr: release-ghcr
	kubectl apply -f config/crd/bases
	kustomize build config/default | kubectl apply -f -
	kubectl apply -f deploy/climatik-operator/manifests/crd.yaml
	kubectl apply -f deploy/climatik-operator/manifests/sample_powercapping.yaml
	kubectl apply -f deploy/climatik-operator/manifests/deployment.yaml

cluster-up: ## setup a cluster for local development
	CLUSTER_PROVIDER=$(CLUSTER_PROVIDER) \
	VERSION=$(LOCAL_DEV_CLUSTER_VERSION) \
	KIND_WORKER_NODES=$(KIND_WORKER_NODES) \
	./hack/cluster.sh up


.PHONY: build
build:
	go build -o bin/manager ./cmd/...

.PHONY: run
run:
	go run ./cmd/...