#!/bin/bash
set -e

CONTAINER_RUNNER="docker"
IMAGE="docker.io/library/golang:bullseye"

$CONTAINER_RUNNER run \
    --rm \
    -it \
    -v `pwd`:/app \
    -p 3333:3333 \
    --pull missing \
    --entrypoint /bin/bash \
    -w /app \
    $IMAGE
