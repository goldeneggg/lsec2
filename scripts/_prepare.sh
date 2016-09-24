#!/bin/sh
PACKAGE=lsec2
PACKAGE_FULL=github.com/goldeneggg/${PACKAGE}

GREP=$(which grep)
SED=$(which sed)
AWK=$(which awk)
SHASUM=$(which shasum)

RELEASES_DIR=$(pwd)/releases
PKG_DIR=$(pwd)/pkg

VERSION=$(${GREP} "const VERSION" version.go | ${SED} -e 's/const VERSION = //g' | ${SED} -e 's/\"//g')
TAG="v${VERSION}"
echo "release tag: ${TAG}"

DEFAULT_OS=('linux' 'darwin' 'windows')
DEFAULT_ARCH=('amd64' '386')
