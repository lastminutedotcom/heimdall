# default for goappfw
APPNAME?="heimdall"
VERSION=$(shell cat .version)

test:
	go test -v ./...

clean:
	go clean -i ./cmd

build: clean
	GOOS=linux GOARCH=amd64 go build  -a -ldflags '-extldflags "-static"' -o $(APPNAME) .

package: build
	docker build --build-arg APP_NAME=$(APPNAME) --build-arg TAG_NAME=stable \
	 	-t registry.bravofly.intra:5000/application/heimdall:$(VERSION) .
	docker push registry.bravofly.intra:5000/application/heimdall:$(VERSION)
	go clean .

.PHONY: test clean build package
