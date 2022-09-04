PKG = $(shell cat go.mod | grep "^module " | sed -e "s/module //g")
NAME = $(shell basename $(PKG))
VERSION = $(shell cat helmx.project.yml|grep version|awk -F : '{print $$2}'|tr -d " ")
COMMIT_SHA ?= $(shell git rev-parse --short HEAD)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CGO_ENABLED ?= 0

GOBUILD=CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -ldflags "-X ${PKG}/version.Version=${VERSION}+sha.${COMMIT_SHA}"
PLATFORM := linux/amd64,linux/arm64
Github_UserName ?= 
Github_Token ?=

WORKSPACE ?= name

clean:
	rm -rf ./cmd/$(WORKSPACE)/out

upgrade:
	go get -u ./...

tidy:
	go mod tidy

build.cm.server: 
	cd ./cmd/$(WORKSPACE) && $(GOBUILD)


build: tidy
	$(MAKE) build.cm tar.cm GOOS=linux GOARCH=amd64
	$(MAKE) build.cm tar.cm GOOS=linux GOARCH=arm64


build.cm:
	cd ./cmd/$(WORKSPACE) && $(GOBUILD) -o ./out/cm-$(GOOS)-$(GOARCH)

tar.cm:
	cd ./cmd/$(WORKSPACE) && tar -czf ./out/cm-$(GOOS)-$(GOARCH).tar.gz -C ./out/ cm-$(GOOS)-$(GOARCH)

install: build.cm
	mv ./cmd/$(WORKSPACE)/out/cm-$(GOOS)-$(GOARCH) ${GOPATH}/bin/cm

docker.client:
	docker buildx build --push --progress plain --platform=${PLATFORM}	\
		--cache-from "type=local,src=/tmp/.buildx-cache" \
		--cache-to "type=local,dest=/tmp/.buildx-cache" \
		--file=./cmd/client/Dockerfile \
		--tag=bryantrh/cm:${VERSION}-${COMMIT_SHA} \
		--build-arg=Github_UserName=${Github_UserName}	\
		--build-arg=Github_Token=${Github_Token}	\
		.

docker.server:
	docker buildx build --push --progress plain --platform=${PLATFORM}	\
		--cache-from "type=local,src=/tmp/.buildx-cache" \
		--cache-to "type=local,dest=/tmp/.buildx-cache" \
		--file=./cmd/server/Dockerfile \
		--tag=bryantrh/cm-server:${VERSION}-${COMMIT_SHA} \
		--build-arg=Github_UserName=${Github_UserName}	\
		--build-arg=Github_Token=${Github_Token}	\
		.


tidy:
	go mod tidy

gen-openapi:
	swag init --pd -d ./cmd/server -o ./cmd/server/docs

gen-client:
	swagger generate client -f ./cmd/server/docs/swagger.json -t ./cmd/client

gen-web:
	npx create-react-app web --template typescript


gen-web-client:
	restful-react import --file ./cmd/server/docs/swagger.json  --output ./cmd/web/src/client-bff.ts