#!/bin/bash
set -eu

source scripts/_prepare.sh

CUR_DATE=$(date "+%Y-%m-%d %H:%M:%S")
COMMIT=$(git log --format=%H -n1)
GO_VERSION=$(go version)

LF="-w -s -extldflags \"-static\""
LF="${LF} -X \"${FULL_PACKAGE}/cmd/lsec2/cli.BuildDate=${CUR_DATE}\""
LF="${LF} -X \"${FULL_PACKAGE}/cmd/lsec2/cli.BuildCommit=${COMMIT}\""
LF="${LF} -X \"${FULL_PACKAGE}/cmd/lsec2/cli.GoVersion=${GO_VERSION}\""

echo "${LF}"
