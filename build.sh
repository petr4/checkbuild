#!/usr/bin/env bash

## Variables must be defined in pipeline
## Uncomment for testing
#
APP_VERSION=v1000

[ -z "$APP_VERSION" ] && echo "APP_VERSION is not set" && exit 1

docker build . -t docker.example.com/checkbuild:${APP_VERSION} --build-arg APP_VERSION=${APP_VERSION}
