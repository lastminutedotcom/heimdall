APPNAME?="heimdall"
VERSION=$(shell cat .version)
REGISTRY?="lmnsre"

test:
	go test -v ./...

clean:
	go clean -i ./cmd

build: clean test
	go build .

docker:
	docker build --build-arg APPNAME=$(APPNAME) -t $(REGISTRY)/$(APPNAME):$(VERSION) .
	docker push $(REGISTRY)/$(APPNAME):$(VERSION)
	rm -f $(APPNAME)

.PHONY: test clean build package