#!/bin/bash

set -eo pipefail

SOURCE_VERSION=${1:-master}
PUBLISH_VERSION=${2:-$SOURCE_VERSION}
ARCH=${3:-amd64}
DOCKER_REPO="bartlettc/wx200"

docker tag ${DOCKER_REPO}:${SOURCE_VERSION}-linux-${ARCH} ${DOCKER_REPO}:${PUBLISH_VERSION}-${ARCH}
echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin
docker push ${DOCKER_REPO}:${PUBLISH_VERSION}-${ARCH}

# On arm build (last build), publish the final combined image
if [[ "${ARCH}" == "arm" ]]; then
# Wait for amd64 image to finish
    docker manifest create \
    ${DOCKER_REPO}:${PUBLISH_VERSION} \
    --amend ${DOCKER_REPO}:${PUBLISH_VERSION}-amd64 \
    --amend ${DOCKER_REPO}:${PUBLISH_VERSION}-arm \
    docker manifest push ${DOCKER_REPO}:${PUBLISH_VERSION}
fi