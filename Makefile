build-image:
	docker build -t quay.io/climatik-project/climatik-operator:latest -f src/climatik-operator/Dockerfile .
push-image: build-image	
	docker push quay.io/climatik-project/climatik-operator:latest


all: build-image