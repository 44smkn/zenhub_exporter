DOCKER_ARCHS ?= amd64 arm64
DOCKER_IMAGE_NAME ?= zenhub-exporter

all:: vet checkmetrics common-all

include Makefile.common

PROMTOOL_DOCKER_IMAGE ?= $(shell docker pull -q quay.io/prometheus/prometheus:latest || echo quay.io/prometheus/prometheus:latest)
PROMTOOL ?= docker run -i --rm -w "$(PWD)" -v "$(PWD):$(PWD)" --entrypoint promtool $(PROMTOOL_DOCKER_IMAGE)

.PHONY: checkmetrics
checkmetrics:
	@echo ">> checking metrics for correctness"
	for file in test/*.metrics; do $(PROMTOOL) check metrics < $$file || exit 1; done