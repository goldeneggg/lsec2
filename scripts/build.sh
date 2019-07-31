#!/bin/bash

set -e

source scripts/_prepare.sh

[ $# -lt 1 ] && { echo 'need to assign output target'; exit 1; }
OUTPUT=${1}
shift

OTHER_OPTS=""
if [ $# -ge 1 ]
then
  OTHER_OPTS="$@"
  echo "OTHER_OPTS = ${OTHER_OPTS}"
fi

LDFLAGS=$(${MYDIR}/_ldflags.sh)

echo "LDFLAGS=${LDFLAGS}, GOOS=${GOOS}, GOARCH=${GOARCH}, OTHER_OPTS=${OTHER_OPTS}"
go build -a -tags netgo -installsuffix netgo -ldflags="${LDFLAGS}" ${OTHER_OPTS} -o ${OUTPUT} ${CMD_GO_DIR}
