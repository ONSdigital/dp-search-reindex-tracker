#!/bin/bash -eux

pushd dp-search-reindex-tracker
  make build
  cp build/dp-search-reindex-tracker Dockerfile.concourse ../build
popd
