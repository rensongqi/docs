FROM alpine:latest AS downloader
WORKDIR /opt
RUN apk add --no-cache curl tar && \
    curl -L https://github.com/tuna/tunasync/releases/download/v0.8.0/tunasync-linux-amd64-bin.tar.gz | tar xz

FROM harbor.rsq.cn/library/ubuntu:18.04
COPY --from=downloader /opt/tunasync /usr/bin/tunasync
COPY --from=downloader /opt/tunasynctl /usr/bin/tunasynctl
ADD start.sh /data/tunasync/
COPY sources.list /etc/apt/sources.list
RUN apt update && apt install rsync -y
ENTRYPOINT ["/data/tunasync/start.sh"]
