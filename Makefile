EXECUTABLE=randrctl
BINDIR=bin
BINARY=${BINDIR}/${EXECUTABLE}
VERSION=$(shell git describe 2>/dev/null || echo "0.0.0")
RELEASE=release
RELEASE_FORMAT=${BINARY}-${GOOS}-${GOARCH}-${VERSION}.tar.gz

default: bin
.PHONY: bin
bin: deps test-deps test ${BINARY}

.PHONY: release
release: release/${EXECUTABLE}-linux-amd64-${VERSION}.tar.gz
release: release/${EXECUTABLE}-linux-386-${VERSION}.tar.gz

.PHONY: release-dir
release-dir:
	mkdir -p ${RELEASE}

${RELEASE}/%-${VERSION}.tar.gz: | release-dir
	GOOS=$(shell echo $* | cut -d '-' -f 2) GOARCH=$(shell echo $* | cut -d '-' -f 3) \
	go build -ldflags="-s -w -X main.version=${VERSION} -v" -o ${RELEASE}/${EXECUTABLE}
	tar -C ${RELEASE} -czf $@ ${EXECUTABLE}
	rm ${RELEASE}/${EXECUTABLE}

.PHONY: bin-dir
bin-dir:
	mkdir -p ${BINDIR}

${BINARY}: bin-dir
	go build -ldflags="-s -w -X main.version=${VERSION} -v" -o $@

.PHONY: lint
lint:
	go get -u github.com/golang/lint/golint
	golint -set_exit_status ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: deps
deps:
	go get ./...

.PHONY: test-deps
test-deps:
	go get -t ./...

.PHONY: clean
clean:
	rm -rf ${BINDIR} 
	rm -rf ${RELEASE}
