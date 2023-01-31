#!/bin/bash
set -e

name=$(basename $0)
root="$(dirname $(dirname $(dirname $0)))"
gomtree=$(go run ${root}/test/realpath/main.go ${root}/gomtree)
left=$(mktemp -d -t go-mtree.XXXXXX)
right=$(mktemp -d -t go-mtree.XXXXXX)

echo "[${name}] Running in ${left} and ${right}"

touch ${left}/one
touch ${left}/two
cp -a ${left}/one ${right}/

$gomtree -K "sha256digest" -p ${left} -c > /tmp/left.mtree
$gomtree -k "sha256digest" -p ${right} -f /tmp/left.mtree
rm -fr ${left} ${right}
