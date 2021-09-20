ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest
LABEL maintainer="44smkn"

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/zenhub_exporter /bin/zenhub_exporter

EXPOSE      9861
USER        nobody
ENTRYPOINT  [ "/bin/zenhub_exporter" ]