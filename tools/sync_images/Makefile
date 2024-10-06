.PHONY: backend

BUILD_TIME=$(shell date +%F-%Z-%H%M%S)

backend:
	- git pull;docker build -t harbor.rsq.cn/projects/sync_docker_images:$(BUILD_TIME) .;docker push harbor.rsq.cn/projects/sync_docker_images:$(BUILD_TIME)