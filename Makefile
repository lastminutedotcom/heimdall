APPNAME?="heimdall"
VERSION=$(shell cat .version)
REGISTRY?="lmnsre"

test:
	go test -v ./...

clean:
	go clean -i ./cmd

build: clean test
	GOOS=linux GOARCH=amd64 go build  -a -ldflags '-extldflags "-static"' -o $(APPNAME) .

package: build
	docker build -t $(REGISTRY)/$(APPNAME):$(VERSION) .
	docker push $(REGISTRY)/$(APPNAME):$(VERSION)
	rm -f $(APPNAME)

.PHONY: test clean build package