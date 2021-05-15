#!/bin/bash -e

# https://hub.docker.com/_/golang
# https://hub.docker.com/_/alpine

if [ "$0" != "./build.sh" ]; then
    echo "Start the build script from the docker folder: ./build.sh" >&2
    exit 1
fi

echo "* Building alpine vigixporter binary"
docker run --rm -v "$PWD/..":/usr/src/github.com/hekmon/vigixporter -w /usr/src/github.com/hekmon/vigixporter golang:1.16.4-alpine3.13 go build -v -ldflags "-s -w" -o docker/vigixporter_alpine
echo
echo "* Building alpine container image"
docker build -t hekmon/vigixporter:1.0.0 -t hekmon/vigixporter:latest .
echo