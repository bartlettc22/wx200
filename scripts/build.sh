#!/bin/bash

set -eo pipefail
# set -o xtrace

VERSION=${1:-master}
PLATFORMS=${2:-"linux/amd64"}
PUBLISH=${3:-"false"}
DOCKER_REPO="bartlettc/wx200"

# archives the binary in the bin/ directory
function archive {
    VERSION=${1}
    OS=${2}
    ARCH=${3}
    cd bin
    if [ "${OS}" == "windows" ]; then
        mv wx200 wx200.exe
        zip wx200-${VERSION}-${OS}-${ARCH}.zip wx200.exe
        rm -rf wx200.exe    
    else 
        tar -cvzf wx200-${VERSION}-${OS}-${ARCH}.tar.gz wx200
        rm -rf wx200
    fi
    cd -
}

if [ "${PUBLISH}" == "true" ]; then
    echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin
fi

# docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
# docker buildx create --name xbuilder --use

LINUX_PLATFORMS=()
OTHER_PLATFORMS=()
for PLATFORM in $(echo ${PLATFORMS} | sed "s/,/ /g")
do
    if [ "${PLATFORM%/*}" == "linux" ]; then
        LINUX_PLATFORMS+=("${PLATFORM}")
    else 
        OTHER_PLATFORMS+=("${PLATFORM}")
    fi
done

mkdir -p bin
rm -rf bin/*

# Build and extract binaries
if [ ${#LINUX_PLATFORMS[@]} > 0 ]; then

    # Build and push linux multiarch images
    PUSH=""
    if "${PUBLISH}" == "true" ]; then
      PUSH="--push"
    fi
    docker buildx build \
        --build-arg VERSION=${VERSION} \
        --platform=$(IFS=, ; echo "${LINUX_PLATFORMS[*]}") \
        --progress plain \
        ${PUSH} \
        -t ${DOCKER_REPO}:${VERSION} \
        .

    # Extract the binaries from the containers
    if "${PUBLISH}" == "true" ]; then 
        for PLATFORM in ${LINUX_PLATFORMS[@]}
        do
            OS=${PLATFORM%/*}
            ARCH=${PLATFORM#*/}
            echo "${OS}/${ARCH}"
            # ARCH=${PLATFORM#"/"}
            docker pull --platform ${PLATFORM} ${DOCKER_REPO}:${VERSION}
            docker run \
                --rm -t \
                --entrypoint "" \
                --name wx200-build-${OS}-${ARCH} \
                -v $(pwd)/bin:/wx200-bin ${DOCKER_REPO}:${VERSION} \
                sh -c "chmod 777 /usr/bin/wx200 && cp /usr/bin/wx200 /wx200-bin"

            # Archive the binary
            archive ${VERSION} ${OS} ${ARCH}
        done
    fi
fi

if [ ${#OTHER_PLATFORMS[@]} > 0 ]; then

    
    for PLATFORM in $(echo ${OTHER_PLATFORMS[@]} | sed "s/,/ /g")
    do
        OS=${PLATFORM%/*}
        ARCH=${PLATFORM#*/}

        # Use docker to build the windows binary
        docker build \
        --build-arg VERSION=${VERSION} \
        --build-arg GOOS=${OS} \
        --build-arg GOARCH=${ARCH} \
        -t ${DOCKER_REPO}:build-${OS}-${ARCH} \
        .

        if "${PUBLISH}" == "true" ]; then 

            # Run the container to extract the binary
            docker run \
            --rm \
            --entrypoint "" \
            --name wx200-build-${OS}-${ARCH} \
            -v $(pwd)/bin:/wx200-bin ${DOCKER_REPO}:build-${OS}-${ARCH} \
            sh -c "cp /usr/bin/wx200 /wx200-bin"

            # Archive the binary
            archive ${VERSION} ${OS} ${ARCH}
        fi
    done
fi