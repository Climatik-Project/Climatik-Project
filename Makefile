.PHONY: tests build-image build-image-ghcr push-image push-image-ghcr release release-ghcr deploy deploy-ghcr cluster-up build run modify-manager-yaml deploy-config deploy-config-ghcr

# Set variables
IMG ?= quay.io/climatik-project/climatik-operator
GITHUB_USERNAME ?= your-github-username
GITHUB_REPO ?= climatik-project
GHCR_IMG ?= ghcr.io/$(GITHUB_USERNAME)/$(GITHUB_REPO)

CLUSTER_PROVIDER ?= kind
LOCAL_DEV_CLUSTER_VERSION ?= main
KIND_WORKER_NODES ?= 2

.DEFAULT_GOAL := default

default: tests build-image-ghcr push-image-ghcr modify-manager-yaml deploy-ghcr

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
	file=$$(cat "deploy/climatik-operator/manifests/deployment.yaml" | sed "s/\$${GITHUB_USERNAME}/$(GITHUB_USERNAME)/g" | sed "s/\$${GITHUB_REPO}/$(GITHUB_REPO)/g"); \
	echo "$$file"; \
	echo "$$file" | kubectl apply -f -

deploy-ghcr: release-ghcr
	kubectl delete deployment operator-powercapping-controller-manager   -n operator-powercapping-system --ignore-not-found
	kubectl apply -f config/crd/bases
	kustomize build config/default | kubectl apply -f -
	kubectl apply -f deploy/climatik-operator/manifests/crd.yaml
	kubectl apply -f deploy/climatik-operator/manifests/sample_powercapping.yaml
	file=$$(cat "deploy/climatik-operator/manifests/deployment.yaml" | sed "s/\$${GITHUB_USERNAME}/$(GITHUB_USERNAME)/g" | sed "s/\$${GITHUB_REPO}/$(GITHUB_REPO)/g"); \
	echo "$$file"; \
	echo "$$file" | kubectl apply -f -

modify-manager-yaml:
	sed -i.bak "s/\$${GITHUB_USERNAME}/$(GITHUB_USERNAME)/g" config/manager/manager.yaml && rm config/manager/manager.yaml.bak
	sed -i.bak "s/\$${GITHUB_REPO}/$(GITHUB_REPO)/g" config/manager/manager.yaml && rm config/manager/manager.yaml.bak

cluster-up: ## setup a cluster for local development
	CLUSTER_PROVIDER=$(CLUSTER_PROVIDER) \
	VERSION=$(LOCAL_DEV_CLUSTER_VERSION) \
	KIND_WORKER_NODES=$(KIND_WORKER_NODES) \
	./hack/cluster.sh up

build:
	go build -o bin/manager ./cmd/...

run:
	go run ./cmd/...

deploy-config: modify-manager-yaml deploy
deploy-config-ghcr: modify-manager-yaml deploy-ghcr