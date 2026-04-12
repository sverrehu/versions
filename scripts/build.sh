#!/bin/bash

BUILDS_DIR=builds
APP_NAME="$(basename "$(pwd)")"

test -d "${BUILDS_DIR}" || mkdir "${BUILDS_DIR}"
env GOOS=linux GOARCH=amd64 go build -o "${BUILDS_DIR}/${APP_NAME}"-linux-amd64
