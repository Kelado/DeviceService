APP_NAME := svc
DOCKER_IMAGE := $(APP_NAME):latest
DOCKER_CONTAINER := $(APP_NAME)_container
PORT := 8000

all: run

build:
	go build -o ./bin/$(APP_NAME)

run: build  
	./bin/$(APP_NAME)

test: 
	go test -cover ./controllers/ ./repositories/

build-image:
	docker build -t $(DOCKER_IMAGE) .

run-image:
	docker run -d --name $(DOCKER_CONTAINER) -p $(PORT):$(PORT) $(DOCKER_IMAGE)

stop-image:
	docker stop $(DOCKER_CONTAINER)

clean-image: stop-image
	docker rm $(DOCKER_CONTAINER)

rmi: clean-image
	docker rmi $(DOCKER_IMAGE)