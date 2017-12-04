#!/bin/bash

set -eu

_gen_readme(){
  cat <<'EOR' > ${TMP_README}
## Install

```
$ cp RELEASE_BIN YOUR_PATH/RELEASE_BIN
```

## Check shasum

```sh
$ shasum -a 256 YOUR_PATH/RELEASE_BIN

# compare to .shasum256 file
$ diff <(cat RELEASE_BIN.shasum256) <(shasum -a 256 YOUR_PATH/RELEASE_BIN | awk '{print $1}')
```
EOR

sed -i -e "s/PACKAGE/${PACKAGE}/g" ${TMP_README}
sed -i -e "s/RELEASE_BIN/${RELEASE_BIN}/g" ${TMP_README}
}

build(){
  local os=${1}
  local arch=${2}

  local release_dir=${RELEASES_DIR}/${PACKAGE}_${os}_${arch}
  [ -d ${release_dir} ] && rm -fr ${release_dir}
  mkdir -p ${release_dir}

  RELEASE_BIN="${PACKAGE}"
  [ ${os} = "windows" ] && RELEASE_BIN=${RELEASE_BIN}.exe

  echo ">>>>>> building: ${release_dir}/${RELEASE_BIN}"

  CUR_DATE=$(date "+%Y-%m-%d %H:%M:%S")
  COMMIT=$(git log --format=%H -n1)
  LD_FLAGS=$(${BASE_DIR}/scripts/_ldflags.sh)
  echo "LD_FLAGS: ${LD_FLAGS}"

  GOOS=${os} GOARCH=${arch} go build -ldflags "${LD_FLAGS}" -o ${release_dir}/${RELEASE_BIN}

  local release_shasum256=${RELEASE_BIN}.shasum256
  ${SHASUM} -a 256 ${release_dir}/${RELEASE_BIN} | ${AWK} '{print $1}' > ${release_dir}/${release_shasum256}
  echo "shasum 256: $(cat ${release_dir}/${release_shasum256})"

  TMP_README=/tmp/readme-${os}-${arch}-${RELEASE_BIN}.md
  _gen_readme
  cp ${TMP_README} ${release_dir}/README.md
}

all(){
  for os in ${DEFAULT_OS[@]}
  do
    for arch in ${DEFAULT_ARCH[@]}
    do
      build ${os} ${arch}
    done
  done
}

source scripts/_prepare.sh

if [ $# -eq 2 ]
then
  echo ">>> release os: ${1}, arch: ${2}"
  build ${1} ${2}
else
  echo ">>> release all platform"
  all
fi

exit 0
