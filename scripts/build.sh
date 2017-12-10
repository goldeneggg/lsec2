#!/bin/bash

set -e

source scripts/_prepare.sh

[ $# -lt 1 ] && { echo 'need to assign output target'; exit 1; }
OUTPUT=${1}

LDFLAGS=$(${MYDIR}/_ldflags.sh)

echo "LDFLAGS=${LDFLAGS}, GOOS=${GOOS}, GOARCH=${GOARCH}"
go build -a -tags netgo -installsuffix netgo -ldflags="${LDFLAGS}" -o ${OUTPUT}
