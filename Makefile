APPNAME?="heimdall"
VERSION=$(shell cat .version)
REGISTRY?="lmnsre"

test:
	go test -v ./...

clean:
	go clean -i ./...

build: clean test
	go build .

docker: test
	docker build --build-arg APPNAME=$(APPNAME) -t $(REGISTRY)/$(APPNAME):$(VERSION) .
	docker push $(REGISTRY)/$(APPNAME):$(VERSION)

.PHONY: test clean build docker