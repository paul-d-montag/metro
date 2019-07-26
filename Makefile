SOURCES := $(shell find . -name '*.go')
PROJECT := metro

.PHONY: build
build: ${PROJECT}

${PROJECT}: ${SOURCES} go.mod
	go get
	go build -v

.PHONY: install
install: build
	mv ${PROJECT} /usr/local/bin/${PROJECT}
