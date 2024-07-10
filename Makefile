.PHONY: all build-image push-image tests release

all: build-image

build-image: tests
	cd python/climatik_operator && docker build -t quay.io/climatik-project/climatik-operator:latest .

push-image: build-image
	docker push quay.io/climatik-project/climatik-operator:latest

tests:
	cd python && PROMETHEUS_HOST="http://localhost:9090" python -m unittest discover tests

release: push-image