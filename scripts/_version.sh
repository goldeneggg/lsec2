#!/bin/bash
MYDIR=$(cd $(dirname $0) && pwd)
BASE_DIR=${MYDIR}/..

GREP=$(which grep)
SED=$(which sed)

${GREP} "const VERSION" ${BASE_DIR}/version.go | ${SED} -e 's/const VERSION = //g' | ${SED} -e 's/\"//g'
