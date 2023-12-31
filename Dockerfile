ARG GOLANG_VERSION=1.21.3
ARG ALPINE_VERSION=3
ARG GV_DOCKERHUB_PROXY
FROM $GV_DOCKERHUB_PROXY/golang:$GOLANG_VERSION AS build

ARG GV_GITLAB_SERVER_READ_GOMODULES
ARG GV_GITLAB_SERVER_READ_GOMODULES_USERNAME
ARG GV_GITLAB_SERVER_READ_GOMODULES_PASSWORD

RUN echo "machine $GV_GITLAB_SERVER_READ_GOMODULES login $GV_GITLAB_SERVER_READ_GOMODULES_USERNAME password $GV_GITLAB_SERVER_READ_GOMODULES_PASSWORD" > ~/.netrc

COPY src /src
WORKDIR /src

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o schema-migration

FROM alpine:$ALPINE_VERSION

RUN adduser --disabled-password --gecos "" db-migration
USER db-migration
COPY --from=build /src/schema-migration /usr/local/bin
CMD ["schema-migration"]
