FROM docker.io/library/golang:1.18-buster AS build-env

FROM build-env AS builder

WORKDIR /go/src
COPY ./ ./

ARG Github_UserName
ARG Github_Token
# build
RUN go env -w GOPRIVATE="github.com/kunlun-qilian" && \ 
    git config --global url."https://${Github_UserName}:${Github_Token}@github.com".insteadOf  "https://github.com" && \
    make build WORKSPACE=client

# runtime
FROM docker.io/library/nginx:1.19.3-alpine

#EXPOSE 80

ARG PROJECT_NAME
ARG PROJECT_VERSION
ARG CM_SERVER_BASEURL
ARG CM_SERVER_USERNAME
ARG CM_SERVER_PASSWORD
ARG KUBE_TUNNEL_GATEWAY_HOST
ENV GOENV=DEV PROJECT_NAME=${PROJECT_NAME} PROJECT_VERSION=${PROJECT_VERSION}
ENV CM_SERVER_BASEURL=${CM_SERVER_BASEURL} CM_SERVER_USERNAME=${CM_SERVER_USERNAME} CM_SERVER_PASSWORD=${CM_SERVER_PASSWORD} KUBE_TUNNEL_GATEWAY_HOST=${KUBE_TUNNEL_GATEWAY_HOST}

#COPY --from=builder /go/src/cmd/client/client /usr/bin/cm
COPY --from=builder /go/src/cmd/client/out/*.tar.gz /usr/share/nginx/html/

#CMD ["/bin/sh"]
