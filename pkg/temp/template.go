package temp

import (
	"bytes"
	"fmt"
	"path/filepath"
)

func Helmxfile() []byte {
	helmxyml := bytes.NewBuffer(nil)

	_, _ = fmt.Fprintln(helmxyml, `
service:
ports:
  - "80"
readinessProbe:
  action: http://:80
  initialDelaySeconds: 5
  periodSeconds: 5
livenessProbe:
  action: http://:80
  initialDelaySeconds: 5
  periodSeconds: 5

	`)
	return helmxyml.Bytes()
}

func Dockerfile(language, service_name string) []byte {
	dockerfile := bytes.NewBuffer(nil)

	switch {
	case language == "go":
		_, _ = fmt.Fprintln(dockerfile, `
FROM docker.io/library/golang:1.18-buster AS build-env

FROM build-env AS builder

WORKDIR /go/src
COPY ./ ./

# build
RUN make build WORKSPACE=`+service_name+`

# runtime
FROM alpine
COPY --from=builder `+filepath.Join("/go/src/cmd", service_name, service_name)+` `+filepath.Join(`/go/bin`, service_name)+`
`)
		fmt.Fprintf(dockerfile, `
EXPOSE 80
ARG PROJECT_NAME
ARG PROJECT_VERSION
ENV GOENV=DEV PROJECT_NAME=${PROJECT_NAME} PROJECT_VERSION=${PROJECT_VERSION}

WORKDIR /go/bin
ENTRYPOINT ["`+filepath.Join(`/go/bin`, service_name)+`"]
`)

	case language == "java":
		_, _ = fmt.Fprintln(dockerfile, `
# First stage: complete build environment
FROM maven:3.5.0-jdk-8-alpine AS builder
# add pom.xml and source code
ADD ./pom.xml pom.xml
ADD ./src src/

# package jar
RUN mvn clean package

# Second stage: minimal runtime environment
FROM openjdk:8-jre-alpine

WORKDIR /java/bin

# copy jar from the first stage
COPY --from=builder `+filepath.Join("target/", service_name)+`-app-1.0-SNAPSHOT.jar  `+filepath.Join(`/java/bin`, service_name)+`-app-1.0-SNAPSHOT.jar`+`

EXPOSE 80
ARG PROJECT_NAME
ARG PROJECT_VERSION
ENV PROJECT_NAME=${PROJECT_NAME} PROJECT_VERSION=${PROJECT_VERSION}


CMD ["java", "-jar", "`+filepath.Join(`/java/bin`, service_name)+`-app-1.0-SNAPSHOT.jar`+`"]
`)
	case language == "python":
		_, _ = fmt.Fprintln(dockerfile, `
FROM python:3.9-alpine as base
FROM base as builder
COPY requirements.txt /requirements.txt
RUN pip install --user -r /requirements.txt

FROM base
# copy only the dependencies installation from the 1st stage image
COPY --from=builder /root/.local /root/.local
COPY . /app
WORKDIR /app

# update PATH environment variable
ENV PATH=/home/app/.local/bin:$PATH

ENTRYPOINT ["python3", "demo.py"]
	`)

	}
	return dockerfile.Bytes()
}

func Makefile() []byte {
	makefile := bytes.NewBuffer(nil)

	_, _ = fmt.Fprintln(makefile, `
PKG = $(shell cat go.mod | grep "^module " | sed -e "s/module //g")
NAME = $(shell basename $(PKG))
VERSION = v$(cat helmx.project.yml|grep version|awk -F : '{print $2}'|tr -d " ")
COMMIT_SHA ?= $(shell git rev-parse --short HEAD)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CGO_ENABLED ?= 0

GOBUILD=CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -ldflags "-X ${PKG}/version.Version=${VERSION}+sha.${COMMIT_SHA}"	OPENAPI=tools openapi

WORKSPACE ?= name

build: 
	cd ./cmd/$(WORKSPACE) && $(GOBUILD)

download:
	go mod download

	`)

	return makefile.Bytes()
}
