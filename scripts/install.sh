#!/bin/bash

set -ex

source scripts/_prepare.sh

OTHER_OPTS=""
if [ $# -ge 1 ]
then
  OTHER_OPTS="$@"
  echo "OTHER_OPTS = ${OTHER_OPTS}"
fi

LDFLAGS=$(${MYDIR}/_ldflags.sh)

echo "LDFLAGS=${LDFLAGS}, GOOS=${GOOS}, GOARCH=${GOARCH}, OTHER_OPTS=${OTHER_OPTS}"
go install -v -ldflags="${LDFLAGS}" ${OTHER_OPTS} ${CMD_GO_DIR}
