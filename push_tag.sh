#!/usr/bin/env bash

# semantic versioning
version=$(cat VERSION)

git tag --force ${version} && \
git push --force --tags