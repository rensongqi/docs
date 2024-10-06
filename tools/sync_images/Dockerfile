FROM harbor.rsq.cn/go/golang:1.21.4-alpine3.17 AS builder
MAINTAINER "songqi.ren@rsq.cn"
USER root
ENV http_proxy http://proxy.rsq.cn:8080
ENV https_proxy http://proxy.rsq.cn:8080
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn
WORKDIR /opt/sync_docker_images
COPY . .
RUN go build -ldflags="-s -w" -o sync_docker_images .

FROM harbor.rsq.cn/library/curlimages/curl:7.85.0-skopeo
USER root
WORKDIR /opt
COPY --from=builder /opt/sync_docker_images/sync_docker_images .
RUN skopeo login -u harbor -p GphVjjGe0 harbor.rsq.cn
ENTRYPOINT ["/opt/sync_docker_images"]