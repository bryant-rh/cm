
FROM docker.io/library/golang:1.17-buster AS build-env

FROM build-env AS builder

WORKDIR /go/src
COPY ./ ./

# build
RUN make build WORKSPACE=server

# runtime
FROM alpine
COPY --from=builder /go/src/cmd/server/server /go/bin/cm-server

EXPOSE 80
COPY --from=builder /go/src/cmd/server/openapi.json /go/bin/openapi.json

ARG PROJECT_NAME
ARG PROJECT_VERSION
ENV GOENV=DEV PROJECT_NAME=${PROJECT_NAME} PROJECT_VERSION=${PROJECT_VERSION}

WORKDIR /go/bin
ENTRYPOINT ["/go/bin/cm-server"]
