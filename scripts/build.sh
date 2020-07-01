#!/bin/bash

set -eo pipefail

VERSION=${1:-master}
GOOS=${2:-linux}
DOCKER_REPO="bartlettc/wx200"

# Directory to house our binaries
mkdir -p bin

# Build the container
docker build --build-arg VERSION=${VERSION} --build-arg GOOS=${GOOS} -t ${DOCKER_REPO}:${VERSION}-${GOOS} ./

# Extract the binary from the container
docker run --rm --entrypoint "" --name wx200-build -v $(pwd)/bin:/wx200-bin ${DOCKER_REPO}:${VERSION}-${GOOS} sh -c "cp /usr/bin/wx200 /wx200-bin"

# Zip up the binary
cd bin
tar -cvzf wx200-${GOOS}-${VERSION}.tar.gz wx200
cd ..