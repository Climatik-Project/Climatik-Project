.PHONY: all build-image push-image tests release

all: build-image

build-image: tests
	docker build -t quay.io/climatik-project/climatik-operator:latest -f python/climatik-operator/Dockerfile .

push-image: build-image
	docker push quay.io/climatik-project/climatik-operator:latest

tests:
	cd python && PROMETHEUS_HOST="example.com" python -m unittest discover tests

release: push-image