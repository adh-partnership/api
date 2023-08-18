#!/bin/bash
#
# Copyright ADH Partnership
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

set -ex

dir="$(cd "$(dirname "$0")" && pwd)"

while [ $# -gt 0 ]; do
    case "$1" in
    --tag)
        TAG="$2"
        shift 2
        ;;
    --image)
        IMAGE="$2"
        shift 2
        ;;
    --hub)
        HUB="$2"
        shift 2
        ;;
    --push)
        PUSH=1
        shift
        ;;
    --help)
        echo "Usage: docker [--tag <tag>] [--image <image>] [--hub <hub>] [--push]"
        echo "  --tag <tag>      Tag to use for the image (defaults to 'latest')"
        echo "  --image <image>  Docker image to build (defaults to 'api')"
        echo "  --hub <hub>      Docker hub to push to (defaults to 'docker.io/denartcc') [example: docker.io/username]"
        echo "  --push           Push the image to the hub"
        echo "  --help: show this help message and exit"
        exit 0
        ;;
    *)
        echo "Unknown argument: $1"
        exit 1
        ;;
    esac
done

HUB=${HUB:-adhp}
IMAGE=${IMAGE:-api}
TAG=${TAG:-latest}

if [[ ! -z "$HUB" ]]; then
    HUB="$HUB/"
fi

pushd $dir/..

make --no-print-directory -e -f Makefile.core.mk test
go get -v ./...
make --no-print-directory -e -f Makefile.core.mk docs
make --no-print-directory -e -f Makefile.core.mk build
docker build . -t $HUB$IMAGE:$TAG

if [[ ! -z "$PUSH" ]]; then
    docker push $HUB$IMAGE:$TAG
fi