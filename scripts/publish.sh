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

formula(){
  cat << EOF > ${PACKAGE}_formula.rb
require "formula"

class Lsec2 < Formula
  homepage 'https://${PACKAGE_FULL}'
  version '${VERSION}'

  if Hardware::CPU.is_32_bit?
    if OS.linux?
      url 'https://${PACKAGE_FULL}/releases/download/${TAG}/${PACKAGE}_${TAG}_linux_386.zip'
      sha256 '$(${SHASUM} -a 256 ${RELEASES_DIR}/${PACKAGE}_linux_386/${PACKAGE}.shasum256 | awk '{print $1}')'
    else
      url 'https://${PACKAGE_FULL}/releases/download/${TAG}/${PACKAGE}_${TAG}_darwin_386.zip'
      sha256 '$(${SHASUM} -a 256 ${RELEASES_DIR}/${PACKAGE}_darwin_386/${PACKAGE}.shasum256 | awk '{print $1}')'
    end
  else
    if OS.linux?
      url 'https://${PACKAGE_FULL}/releases/download/${TAG}/${PACKAGE}_${TAG}_linux_amd64.zip'
      sha256 '$(${SHASUM} -a 256 ${RELEASES_DIR}/${PACKAGE}_linux_amd64/${PACKAGE}.shasum256 | awk '{print $1}')'
    else
      url 'https://${PACKAGE_FULL}/releases/download/${TAG}/${PACKAGE}_${TAG}_darwin_amd64.zip'
      sha256 '$(${SHASUM} -a 256 ${RELEASES_DIR}/${PACKAGE}_darwin_amd64/${PACKAGE}.shasum256 | awk '{print $1}')'
    end
  end

  def install
    bin.install 'lsec2'
  end
end
EOF
}

source scripts/_prepare.sh && publish && formula
