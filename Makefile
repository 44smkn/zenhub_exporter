EXE =
ifeq ($(GOOS),windows)
EXE = .exe
endif

## The following tasks delegate to `script/build.go` so they can be run cross-platform.

.PHONY: bin/s3select$(EXE)
bin/s3select$(EXE): script/build
	@script/build $@

script/build: script/build.go
	GOOS= GOARCH= GOARM= GOFLAGS= CGO_ENABLED= go build -o $@ $<

.PHONY: clean
clean: script/build
	@script/build $@

# just a convenience task around `go test`
.PHONY: test
test:
	go generate ./...
	go test ./...