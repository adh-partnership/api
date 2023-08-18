#!/bin/bash
#
# Copyright Daniel Hawton. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -x

SCRIPTPATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

OUT=${1:?"output path"}
shift

set -e

export BUILD_GOOS=${GOOS:-linux}
export BUILD_GOARCH=${GOARCH:-amd64}
GOBINARY=${GOBINARY:-go}
GOPKG="$GOPATH/pkg"
BUILDINFO=${BUILDINFO:-""}
STATIC=${STATIC:-1}
LDFLAGS=${LDFLAGS:--extldflags -static}
GOBUILDFLAGS=${GOBUILDFLAGS:-""}
# Split GOBUILDFLAGS by spaces into an array called GOBUILDFLAGS_ARRAY.
IFS=' ' read -r -a GOBUILDFLAGS_ARRAY <<< "$GOBUILDFLAGS"

GCFLAGS=${GCFLAGS:-}
export CGO_ENABLED=${CGO_ENABLED:-0}

if [[ "${STATIC}" !=  "1" ]];then
    LDFLAGS=""
fi

OPTIMIZATION_FLAGS=(-trimpath)
if [ "${DEBUG}" == "1" ]; then
    OPTIMIZATION_FLAGS=()
fi

time GOOS=${BUILD_GOOS} GOARCH=${BUILD_GOARCH} ${GOBINARY} build \
        ${V} "${GOBUILDFLAGS_ARRAY[@]}" ${GCFLAGS:+-gcflags "${GCFLAGS}"} \
        -o "${OUT}" \
        "${OPTIMIZATION_FLAGS[@]}" \
        -pkgdir="${GOPKG}/${BUILD_GOOS}_${BUILD_GOARCH}" \
        -ldflags "${LDFLAGS} ${LD_EXTRAFLAGS}" "${@}"

static_files=(boundaries.json init.sh config.yaml.example)
static_dirs=(static docs)

for file in "${static_files[@]}"; do
  cp $file out/$file
done

for dir in "${static_dirs[@]}"; do
  cp -Rp $dir out/$dir
done