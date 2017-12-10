#!/bin/bash
MYDIR=$(cd $(dirname $0) && pwd)
BASE_DIR=${MYDIR}/..
RELEASES_DIR=${BASE_DIR}/releases
PKG_DIR=${BASE_DIR}/pkg

PACKAGE=lsec2
PACKAGE_FULL=github.com/goldeneggg/${PACKAGE}
FORMULA_CLASS=Lsec2

AWK=$(which awk)
SHASUM=$(which shasum)

VERSION=$(${BASE_DIR}/scripts/_version.sh)
TAG=v${VERSION}

DEFAULT_OS=('linux' 'darwin' 'windows' 'freebsd')
DEFAULT_ARCH=('amd64' '386')
