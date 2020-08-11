#!/usr/bin/env bash

# shellcheck source=./build.conf
source "$(pwd)/build.conf"
git tag -f "${VERSION}"
git push -f nkryuchkov "${VERSION}" || exit 1
