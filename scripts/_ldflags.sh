#!/bin/bash
MYDIR=$(cd $(dirname $0) && pwd)
BASE_DIR=${MYDIR}/..

CUR_DATE=$(date "+%Y-%m-%d %H:%M:%S")
COMMIT=$(git log --format=%H -n1)
LF="-w -s -extldflags \"-static\""
LF="${LF} -X \"main.version=$(${BASE_DIR}/scripts/_version.sh)\""
LF="${LF} -X \"main.buildDate=${CUR_DATE}\""
LF="${LF} -X \"main.buildCommit=${COMMIT}\""

echo "${LF}"
