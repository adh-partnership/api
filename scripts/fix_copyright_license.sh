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

set -e

WD=$(dirname "$0")
WD=$(cd "$WD"; pwd)

for fn in "$@"; do
  if ! grep -L -q -e "Apache License, Version 2" -e "Copyright" "${fn}"; then
    if [[ "${fn}" == *.go ]]; then
      echo "Fixing license: ${fn}"
      newfile=$(cat "${WD}/banners/golang.txt" <(echo) "${fn}")
      echo "${newfile}" > "${fn}"
    elif [[ "${fn}" == *.sh ]]; then
      echo "Fixing license: ${fn}"
      newfile=$(cat "${WD}/banners/shell.txt" <(echo) "${fn}")
      echo "${newfile}" > "${fn}"
    else
      echo "Cannot fix license: ${fn}. Not supported file type."
    fi
  fi
done