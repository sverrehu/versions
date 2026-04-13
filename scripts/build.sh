#!/bin/bash

BUILDS_DIR=builds
APP_NAME="$(basename "$(pwd)")"

test -d "${BUILDS_DIR}" || mkdir "${BUILDS_DIR}"
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -trimpath -ldflags="-s -w" -o "${BUILDS_DIR}/${APP_NAME}"-linux-amd64
