#!/bin/sh

set -eu

publish(){
  [ -d ${PKG_DIR} ] && rm -r ${PKG_DIR}
  mkdir -p ${PKG_DIR}

  pushd ${RELEASES_DIR}
  for os in ${DEFAULT_OS[@]}
  do
    for arch in ${DEFAULT_ARCH[@]}
    do
      local build_name=${PACKAGE}_${os}_${arch}
      zip -9 ${PKG_DIR}/${build_name}.zip ${build_name}/*
      echo "created artifacts ${PKG_DIR}/${build_name}.zip"
    done
  done
  popd
  ghr --draft --replace ${TAG} ${PKG_DIR}
}

shasum256() {
  local os=${1}
  local arch=${2}

  ${SHASUM} -a 256 ${PKG_DIR}/${PACKAGE}_${os}_${arch}.zip | awk '{print $1}'
}

formula(){
  cat << EOF > ${PACKAGE}_formula.rb
require "formula"

class Lsec2 < Formula
  homepage 'https://${PACKAGE_FULL}'
  version '${VERSION}'

  if Hardware::CPU.is_32_bit?
    if OS.linux?
      url 'https://${PACKAGE_FULL}/releases/download/${TAG}/${PACKAGE}_linux_386.zip'
      sha256 '$(shasum256 linux 386)'
    else
      url 'https://${PACKAGE_FULL}/releases/download/${TAG}/${PACKAGE}_darwin_386.zip'
      sha256 '$(shasum256 darwin 386)'
    end
  else
    if OS.linux?
      url 'https://${PACKAGE_FULL}/releases/download/${TAG}/${PACKAGE}_linux_amd64.zip'
      sha256 '$(shasum256 linux amd64)'
    else
      url 'https://${PACKAGE_FULL}/releases/download/${TAG}/${PACKAGE}_darwin_amd64.zip'
      sha256 '$(shasum256 darwin amd64)'
    end
  end

  def install
    bin.install '${PACKAGE}'
  end
end
EOF
}

source scripts/_prepare.sh

if [ $# -eq 1 ]
then
  [ ${1} = "formula-only" ] && formula
else
  publish && formula
fi

