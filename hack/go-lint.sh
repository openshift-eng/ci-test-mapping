#!/bin/bash
# If CI is true, run golangci-lint directly. If we're on a developer's
# local machine, run golangci-lint from a container so we're ensuring
# a consistent environment.

set -ex


export GOMEMLIMIT=200MiB  # lint can bloat memory, make it stay slim
if [ "$CI" = "true" ];
then
  # lint will balloon if allowed; keep it limited
  go version
  golangci-lint version -v
  golangci-lint --timeout 10m "${@}"
else
  DOCKER=${DOCKER:-podman}

  if ! which "$DOCKER" > /dev/null 2>&1;
  then
    echo "$DOCKER not found, please install."
    exit 1
  fi

  $DOCKER run --rm \
    --volume "${PWD}:/go/src/github.com/openshift-eng/ci-test-mapping:z" \
    --workdir /go/src/github.com/openshift-eng/ci-test-mapping \
    --env GOMEMLIMIT \
    docker.io/golangci/golangci-lint:v2.10.1 \
    golangci-lint "${@}"
fi
