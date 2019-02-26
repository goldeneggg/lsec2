#!/bin/bash

set -ex

source scripts/_prepare.sh

[ $# -lt 1 ] && { echo 'need to assign output target'; exit 1; }
OUTPUT=${1}
shift

if [ $# -ge 1 ]
then
  OTHER_OPTS="$@"
  echo "OTHER_OPTS = ${OTHER_OPTS}"
fi

LDFLAGS=$(${MYDIR}/_ldflags.sh)

echo "LDFLAGS=${LDFLAGS}, GOOS=${GOOS}, GOARCH=${GOARCH}"
go build -a -tags netgo -installsuffix netgo -ldflags="${LDFLAGS}" ${OTHER_OPTS} -o ${OUTPUT}
