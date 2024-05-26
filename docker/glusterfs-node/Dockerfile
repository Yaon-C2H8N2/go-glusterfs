FROM debian:latest

WORKDIR /go-glusterfs

RUN apt-get update && apt-get install -y glusterfs-server golang

ENTRYPOINT ["glusterd", "-N"]