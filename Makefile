.PHONY: all tests build-image-ghcr push-image-ghcr release-ghcr deploy-ghcr

# Set variables
IMG ?= quay.io/climatik-project/climatik-operator
GITHUB_USERNAME ?= your-github-username
GITHUB_REPO ?= climatik-project
GHCR_IMG ?= ghcr.io/$(GITHUB_USERNAME)/$(GITHUB_REPO)

all: deploy-ghcr

tests: build-image-ghcr
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

deploy-ghcr: release-ghcr
	kubectl apply -f config/crd/bases
	kustomize build config/default | kubectl apply -f -