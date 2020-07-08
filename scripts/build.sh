#!/bin/bash

set -eo pipefail
set -o xtrace

VERSION=${1:-master}
OS=${2:-linux}
ARCH=${3:-amd64}
DOCKER_REPO="bartlettc/wx200"


# docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
# docker buildx create --name xbuilder --use
# Don't build an ARM Windows binary
# if [[ "${OS}" == "windows" && "${ARCH}" == "arm" ]]; then
#     exit 0
# fi

# Directory to house our binaries
mkdir -p bin

# Build the container
# docker build \
#     --build-arg VERSION=${VERSION} \
#     --build-arg GOOS=${OS} \
#     --build-arg GOARCH=${ARCH} \
#     -t ${DOCKER_REPO}:${VERSION}-${OS}-${ARCH} \
#     ./
docker buildx build --progress plain --output=type=local,dest=./bin/test --platform=linux/arm64 -t test .
# docker buildx -h
# docker buildx build \
#      --progress plain \
#     --platform=linux/amd64,linux/arm64 \
#     -t ${DOCKER_REPO}:${VERSION} \
#     .

docker image ls -a
ls -al ./bin/test

docker run test
# docker manifest create bartlettc/wx200:test \
# --amend bartlettc/wx200:test-amd64 \
# --amend bartlettc/wx200:test-arm

# Extract the binary from the container
# docker run \
#     --rm \
#     --entrypoint "" \
#     --name wx200-build \
#     -v $(pwd)/bin:/wx200-bin ${DOCKER_REPO}:${VERSION}-${OS} \
#     sh -c "cp /usr/bin/wx200 /wx200-bin"

# # Zip up the binary
# cd bin
# if [[ "${OS}" == "linux" ]]; then
#     tar -cvzf wx200-${VERSION}-${OS}-${ARCH}.tar.gz wx200
#     rm wx200
# elif [[ "${OS}" == "windows" ]]; then
#     mv wx200 wx200.exe
#     zip wx200-${VERSION}-${OS}-${ARCH}.zip wx200.exe
#     rm wx200.exe
# fi
# cd ..