#!/bin/bash
set -e

name=$(basename $0)
root="$(dirname $(dirname $(dirname $0)))"
gomtree=$(go run ${root}/test/realpath/main.go ${root}/gomtree)
t=$(mktemp -d -t go-mtree.XXXXXX)

echo "[${name}] Running in ${t}"

pushd ${root}
mkdir -p ${t}/extract
git archive --format=tar HEAD^{tree} . | tar -C ${t}/extract/ -x

${gomtree} -k sha1digest -c -p ${t}/extract/ > ${t}/${name}.mtree
${gomtree} -f ${t}/${name}.mtree -k md5digest -p ${t}/extract/

popd
rm -rf ${t}
