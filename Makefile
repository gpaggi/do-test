PROG = echoapi
VERSION = $$(cat VERSION)
IMG_NAME = gpaggi/echoapi

.PHONY: build build-docker clean distribute clean-image fmt test

build:
	@cd src; CGO_ENABLED=0 go build -v -ldflags "-X github.com/gpaggi/$(PROG)/version.Version=$(VERSION) -s -w" -o ../bin/$(PROG) .

clean:
	@rm bin/$(PROG)

build-docker:
	@cd src; docker build --build-arg VERSION=$(VERSION) -t $(IMG_NAME):$(VERSION) -f ../Dockerfile .
	@cd src; docker tag $(IMG_NAME):$(VERSION) $(IMG_NAME):latest

distribute: build-docker
	@docker push $(IMG_NAME):$(VERSION)
	@docker push $(IMG_NAME):latest

clean-image:
	@docker image rm $(IMG_NAME):$(VERSION)

fmt:
	@cd src; go mod tidy
	@cd src; find . -type f -name '*.go' | xargs gofmt -w -s -e

test:
	@cd src; go test -failfast -v ./...