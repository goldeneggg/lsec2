#!/bin/bash
set -eu

source scripts/_prepare.sh

GREP=$(which grep)
SED=$(which sed)

${GREP} "const VERSION" ${CMD_GO_DIR}/version/version.go | ${SED} -e 's/const VERSION = //g' | ${SED} -e 's/\"//g'
